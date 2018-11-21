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
	ip := "35.184.170.191"
	port := "80"
	address := fmt.Sprintf("%s:%s", ip, port)

	fmt.Printf("Server listening on %s at port %s...\n", ip, port)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		fmt.Printf("Error - failed to listen on %s\n", address)
		fmt.Printf("%s\n", err.Error())
	}
}
