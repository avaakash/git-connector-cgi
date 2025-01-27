package main

import (
	"net/http"

	"github.com/harness/github-connector-cgi/handler"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	http.HandleFunc("/", handler.HandleRequest)
	err := http.ListenAndServe(":8080", nil)
	// err := cgi.Serve(http.DefaultServeMux)

	if err != nil {
		logrus.WithError(err).Fatal("Failed to serve CGI")
	}
}
