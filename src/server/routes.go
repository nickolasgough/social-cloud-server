package server

import (
	userApi "server/src/internal/user/api"
	"server/src/server/endpoint"
)

func (s *Server) Routes() map[string]endpoint.Handler {
	return map[string]endpoint.Handler{
		"/user/create": userApi.NewCreateHandler(),
	}
}
