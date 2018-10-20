package main

import (
	"net/http"

	"server/src/server"
)

func main() {
	server := server.NewServer(http.DefaultServeMux)
	server.RegisterRoutes()
	server.ListenAndServe()
}
