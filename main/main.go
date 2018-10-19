package main

import (
	"net/http"

	"cloud-receipts/src/server"
)

func main() {
	server := server.NewServer(http.DefaultServeMux)
	server.RegisterRoutes()
	server.ListenAndServe()
}
