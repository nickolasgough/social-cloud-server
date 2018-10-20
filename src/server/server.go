package server

import (
	"net/http"
	"fmt"

	"social-cloud-server/src/server/endpoint"
)

type Server struct {
	Client *http.ServeMux
}

func NewServer(client *http.ServeMux) *Server {
	return &Server{
		Client: client,
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
