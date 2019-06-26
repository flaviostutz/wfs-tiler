package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/boundlessgeo/wfs3/model"
	"github.com/boundlessgeo/wfs3/ogc"
	"github.com/gin-gonic/gin"
)

func (h *HTTPServer) setupVectorTilerHandlers(wfsURL string) {
	h.router.GET("/{collection}/{z}/{y}/{x}.mvt", getVectorTile(wfsURL))
}

func getVectorTile(wfsURL string) func(*gin.Context) {
	return func(c *gin.Context) {
		https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames#Implementations
		https://github.com/murphy214/vector-tile-go
		"application/vnd.mapbox-vector-tile"
		// c.JSON(http.StatusOK, gin.H{"collections": cis})
	}
}

