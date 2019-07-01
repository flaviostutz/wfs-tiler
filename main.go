package main

import (
	"flag"
	"os"

	"github.com/flaviostutz/wfs-tiler/handlers"
	"github.com/sirupsen/logrus"
)

func main() {
	logLevel := flag.String("loglevel", "debug", "debug, info, warning, error")
	wfsURL := flag.String("wfs-url", "", "WFS 3.0 server API URL from which to get feature in order to provide the vector tile contents")
	cacheControl := flag.String("cache-control", "", "HTTP response Cache-Control header contents for all requests. If empty, no header is set.")
	simplificationLevel := flag.Int("simplify-level", 10, "Geometry simplification level on min (wider) zoom level. The simplificication is decreased to 0 as zoom level approaches min (more detailed)")
	minGeomLength := flag.Int("min-geom-length", 3600, "At min zoom (wider), geometries with length less than this value are hidden. This parameter is decreased to 0 as zoom level approaches min (more detailed)")
	maxZoomLevel := flag.Int("max-zoom-level", 18, "Max allowed zoom level")
	flag.Parse()

	switch *logLevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
		break
	case "warning":
		logrus.SetLevel(logrus.WarnLevel)
		break
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
		break
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	opt := handlers.Options{
		WFSURL:              *wfsURL,
		CacheControl:        *cacheControl,
		SimplificationLevel: *simplificationLevel,
		MinGeomLength:       *minGeomLength,
		MaxZoomLevel:        *maxZoomLevel,
	}

	if opt.WFSURL == "" {
		logrus.Errorf("'--wfs-url' is required")
		os.Exit(1)
	}

	logrus.Infof("====Starting WFS-TILER====")
	h := handlers.NewHTTPServer(opt)
	err := h.Start()
	if err != nil {
		logrus.Errorf("Error starting server. err=%s", err)
		os.Exit(1)
	}

}
