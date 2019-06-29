package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"

	// vt "github.com/murphy214/vector-tile-go"
	"github.com/paulmach/orb/encoding/mvt"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/maptile"
	"github.com/paulmach/orb/simplify"
	"github.com/sirupsen/logrus"
)

func (h *HTTPServer) setupVectorTilerHandlers(wfsURL string, cacheControl string) {
	h.router.GET("/tiles/:collection/:z/:x/:y.mvt", getVectorTile(wfsURL, cacheControl))
}

func getVectorTile(wfsURL string, cacheControl string) func(*gin.Context) {
	return func(c *gin.Context) {
		collection := c.Param("collection")
		x, err := strconv.Atoi(c.Param("x"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "'x' parameter is invalid"})
			return
		}
		vy := c.Param("y.mvt")
		vy = strings.Split(vy, ".")[0]
		y, err := strconv.Atoi(vy)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "'y' parameter is invalid"})
			return
		}
		z, err := strconv.Atoi(c.Param("z"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "'z' parameter is invalid"})
			return
		}

		bbox := tile2BBox(x, y, z)
		bboxstr := fmt.Sprintf("%f,%f,%f,%f", bbox.west, bbox.south, bbox.east, bbox.north)
		logrus.Debugf("tile=%d,%d,%d; bbox=%s", x, y, z, bboxstr)

		limitstr := c.Query("limit")
		if limitstr != "" {
			limitstr = fmt.Sprintf("&limit=%s", limitstr)
		}

		timestr := c.Query("time")
		if timestr != "" {
			timestr = fmt.Sprintf("&time=%s", timestr)
		}

		propertiesFilterStr := ""
		params := c.Request.URL.Query()
		for k, v := range params {
			if k != "time" && k != "bbox" && k != "limit" {
				propertiesFilterStr = fmt.Sprintf("%s&%s=%s", propertiesFilterStr, k, v)
			}
		}

		if timestr != "" {
			timestr = fmt.Sprintf("&time=%s", timestr)
		}

		q := fmt.Sprintf("%s/collections/%s/items?bbox=%s%s%s%s", wfsURL, collection, bboxstr, limitstr, timestr, propertiesFilterStr)
		logrus.Debugf("WFS query: %s", q)
		resp, err := http.Get(q)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Error requesting WFS service. err=%s", err)})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			msg := fmt.Sprintf("WFS invocation status != 200. status=%d", resp.StatusCode)
			c.JSON(resp.StatusCode, gin.H{"message": msg})
			return
		}

		var fc geojson.FeatureCollection
		data, err0 := ioutil.ReadAll(resp.Body)
		if err0 != nil {
			msg := fmt.Sprintf("Error reading WFS service response. err=%s", err0)
			logrus.Errorf(msg)
			c.JSON(http.StatusInternalServerError, gin.H{"message": msg})
			return
		}

		logrus.Debugf("WFS response bytes: %d", len(data))
		err = json.Unmarshal(data, &fc)
		if err != nil {
			msg := fmt.Sprintf("Error parsing WFS service response. err=%s", err)
			logrus.Errorf(msg)
			c.JSON(http.StatusInternalServerError, gin.H{"message": msg})
			return
		}
		logrus.Debugf("WFS response feature count: %d", len(fc.Features))

		collections := make(map[string]*geojson.FeatureCollection)
		collections[collection] = &fc

		// Convert to a layers object and project to tile coordinates.
		layers := mvt.NewLayers(collections)
		layers.ProjectToTile(maptile.New(uint32(x), uint32(y), maptile.Zoom(z)))

		// In order to be used as source for MapboxGL geometries need to be clipped
		// to max allowed extent. (uncomment next line)
		layers.Clip(mvt.MapboxGLDefaultExtentBound)

		// Simplify the geometry now that it's in tile coordinate space.
		layers.Simplify(simplify.DouglasPeucker(10.0))

		// Depending on use-case remove empty geometry, those too small to be
		// represented in this tile space.
		// In this case lines shorter than 1, and areas smaller than 2.
		layers.RemoveEmpty(30.0*(18.0/float64(z)), 200.0*(18.0/float64(z)))

		// encoding using the Mapbox Vector Tile protobuf encoding.
		layerBytes, err0 := mvt.Marshal(layers)
		if err0 != nil {
			msg := fmt.Sprintf("Error generating MVT bytes. err=%s", err0)
			logrus.Errorf(msg)
			c.JSON(http.StatusInternalServerError, gin.H{"message": msg})
			return
		}

		if cacheControl != ""{
			c.Header("Cache-Control", cacheControl)
		}
		c.Render(
			http.StatusOK, render.Data{
				ContentType: "application/vnd.mapbox-vector-tile",
				Data:        layerBytes,
			})

	}
}
