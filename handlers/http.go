package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/flaviostutz/wfs3-tiler/handlers"
	"github.com/gin-gonic/gin"

	cors "github.com/itsjamie/gin-cors"
)

type HTTPServer struct {
	server *http.Server
	router *gin.Engine
}

func NewHTTPServer() *HTTPServer {
	router := gin.Default()

	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET",
		RequestHeaders:  "Origin, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          24 * 3600 * time.Second,
		Credentials:     false,
		ValidateHeaders: false,
	}))

	setupVectorTilerHandlers(wfsURL)

	return &HTTPServer{server: &http.Server{
		Addr:    ":3000",
		Handler: router,
	}, router: router}
}

//Start the main HTTP Server entry
func (s *HTTPServer) Start(wfsURL string) error {
	logrus.Infof("Starting HTTP Server on port 3000")
	return s.server.ListenAndServe()
}

//Stop for the run group
func (s *HTTPServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	logrus.Warnf("Stopping HTTP Server")
	return s.server.Shutdown(ctx)
}
