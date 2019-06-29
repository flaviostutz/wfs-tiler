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
	m "github.com/murphy214/mercantile"
	vt "github.com/murphy214/vector-tile-go"
	geojson "github.com/paulmach/go.geojson"
	"github.com/sirupsen/logrus"
)

func (h *HTTPServer) setupVectorTilerHandlers(wfsURL string) {
	h.router.GET("/tiles/:collection/:z/:x/:y.mvt", getVectorTile(wfsURL))
}

func getVectorTile(wfsURL string) func(*gin.Context) {
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

		logrus.Debugf("WFS response: %s", data)
		err = json.Unmarshal(data, &fc)
		if err != nil {
			msg := fmt.Sprintf("Error parsing WFS service response. err=%s", err)
			logrus.Errorf(msg)
			c.JSON(http.StatusInternalServerError, gin.H{"message": msg})
			return
		}
		features := fc.Features

		xyz := m.TileID{
			X: int64(x),
			Y: int64(y),
			Z: uint64(z)}

		config1 := vt.NewConfig(collection, xyz)
		layer1bytes, err := vt.WriteLayer(features, config1)
		if err != nil {
			fmt.Println(err)
		}

		// c.Header("Cache-Control", "public, max-age=604800")
		c.Render(
			http.StatusOK, render.Data{
				ContentType: "application/vnd.mapbox-vector-tile",
				Data:        layer1bytes,
			})

	}
}
