package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	cors "github.com/itsjamie/gin-cors"
)

type HTTPServer struct {
	server *http.Server
	router *gin.Engine
}

func NewHTTPServer(wfsURL string) *HTTPServer {
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

	h := &HTTPServer{server: &http.Server{
		Addr:    ":3000",
		Handler: router,
	}, router: router}

	h.setupVectorTilerHandlers(wfsURL)

	return h
}

//Start the main HTTP Server entry
func (s *HTTPServer) Start() error {
	logrus.Infof("Starting HTTP Server on port 3000")
	return s.server.ListenAndServe()
}

//Stop for the run group
// func (s *HTTPServer) Stop() error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	logrus.Warnf("Stopping HTTP Server")
// 	return s.server.Shutdown(ctx)
// }
