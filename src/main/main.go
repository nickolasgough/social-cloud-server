package main

import (
	"net/http"
	"fmt"
	"os"

	"social-cloud-server/src/database"
	"social-cloud-server/src/server"
)

func main() {
	db := database.NewDatabase()
	err := db.Connect()
	if err != nil {
		fmt.Printf("Failed to connect to database: %s\n", err.Error())
		os.Exit(1)
	}

	s := server.NewServer(http.DefaultServeMux, db)
	s.RegisterRoutes()
	s.ListenAndServe()
}
