package server

import (
	userApi "social-cloud-server/src/internal/user/api"
	"social-cloud-server/src/server/endpoint"
)

func (s *Server) Routes() map[string]endpoint.Handler {
	return map[string]endpoint.Handler{
		"/user/create": userApi.NewCreateHandler(),
	}
}
