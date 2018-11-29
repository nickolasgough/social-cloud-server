package server

import (
	"net/http"
	"fmt"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
	"social-cloud-server/src/bucket"
)

type Server struct {
	Client   *http.ServeMux
	Database *database.Database
	Bucket *bucket.Bucket
}

func NewServer(cl *http.ServeMux, db *database.Database, b *bucket.Bucket) *Server {
	return &Server{
		Client: cl,
		Database: db,
		Bucket: b,
	}
}

func (s *Server) RegisterRoutes() error {
	for r, h := range s.Routes() {
		s.RegisterHandler(r, h)
	}

	return nil
}

func (s *Server) RegisterHandler(route string, handler endpoint.Handler) {
	s.Client.Handle(route, &endpoint.Listener{Handler: handler})
}

func (s *Server) ListenAndServe() error {
	port := "8080"
	address := fmt.Sprintf(":%s", port)

	fmt.Printf("Server listening on port %s...\n", port)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		return err
	}

	return nil
}
