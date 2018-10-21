package server

import (
	"net/http"
	"fmt"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
)

type Server struct {
	Client   *http.ServeMux
	Database *database.Database
}

func NewServer(cl *http.ServeMux, db *database.Database) *Server {
	return &Server{
		Client: cl,
		Database: db,
	}
}

func (s *Server) RegisterRoutes() {
	for r, h := range s.Routes() {
		s.RegisterHandler(r, h)
	}
}

func (s *Server) RegisterHandler(route string, handler endpoint.Handler) {
	s.Client.Handle(route, &endpoint.Listener{Handler: handler})
}

func (s *Server) ListenAndServe() {
	fmt.Printf("Server listening on localhost at port 8080...\n")
	http.ListenAndServe(":8080", nil)
}
