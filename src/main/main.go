package main

import (
	"net/http"

	"social-cloud-server/src/server"
)

func main() {
	server := server.NewServer(http.DefaultServeMux)
	server.RegisterRoutes()
	server.ListenAndServe()
}
