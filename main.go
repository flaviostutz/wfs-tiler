package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logLevel := flag.String("loglevel", "debug", "debug, info, warning, error")
	wfsURL0 := flag.Bool("wfs-url", "", "WFS 3.0 server API URL from which to get feature in order to provide the vector tile contents")
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

	logrus.Infof("====Starting WFS-TILER====")

	wfsURL := *wfsURL0
	if wfsURL == "" {
		logrus.Errorf("'wfsUrl' is required")
		exit(1)
	}

	h := handlers.NewHTTPServer(wfsUrl)
	err := h.Start()
	if err != nil {
		logrus.Errorf(err)
		exit(1)
	}

	// var dbErr, httpErr error

	// go func() {
	// 	err := h.Start()
	// 	if err != nil {
	// 		logrus.Errorf(err)
	// 		exit(1)
	// 	}
	// }()

	// sigchan := make(chan os.Signal, 1)
	// signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// running := true
	// for running == true {
	// 	select {
	// 	case sig := <-sigchan:
	// 		h.Stop(httpErr)
	// 		running = false
	// 	}
	// }

}
