package main

import (
	"net/http"
	"net/http/cgi"

	"github.com/harness/git-connector-cgi/handler"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	http.HandleFunc("/", handler.HandleRequest)
	// err := http.ListenAndServe(":9600", nil) // nil uses the default ServeMux
	err := cgi.Serve(http.DefaultServeMux)

	if err != nil {
		logrus.WithError(err).Fatal("Failed to serve CGI")
	}
}
