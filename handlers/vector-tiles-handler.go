package handlers

import (
	"github.com/gin-gonic/gin"
	// "github.com/sirupsen/logrus"
)

func (h *HTTPServer) setupVectorTilerHandlers(wfsURL string) {
	h.router.GET("/tiles/{collection}/{z}/{x}/{y}.mvt", getVectorTile(wfsURL))
}

func getVectorTile(wfsURL string) func(*gin.Context) {
	return func(c *gin.Context) {
		// https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames#Implementations
		// https://github.com/murphy214/vector-tile-go
		// "application/vnd.mapbox-vector-tile"
		// c.JSON(http.StatusOK, gin.H{"collections": cis})
	}
}
