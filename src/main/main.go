package main

import (
	"net/http"
	"fmt"
	"os"
	"context"

	"social-cloud-server/src/database"
	"social-cloud-server/src/server"
	"social-cloud-server/src/bucket"
)

func main() {
	db := database.NewDatabase()
	err := db.ConnectDatabase()
	if err != nil {
		fmt.Printf("Failed to connect to database: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println("Successfully established a connection to the database")

	//err = db.BuildModels()
	//if err != nil {
	//	fmt.Printf("Failed to construct the database: %s\n", err.Error())
	//	os.Exit(1)
	//}
	//fmt.Println("Successfully rebuilt the application models")

	b := bucket.NewBucket()
	err = b.ConnectBucket(context.Background())
	if err != nil {
		fmt.Printf("Failed to connect to bucket: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println("Successfully established a connection to the bucket")

	s := server.NewServer(http.DefaultServeMux, db, b)
	err = s.RegisterRoutes()
	if err != nil {
		fmt.Printf("Failed to register server routes: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println("Successfully registered the server handlers for the app")

	err = s.ListenAndServe()
	if err != nil {
		fmt.Printf("Failed to listen and serve requests: %s\n", err.Error())
		os.Exit(1)
	}
}
