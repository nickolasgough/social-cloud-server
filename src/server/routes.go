package server

import (
	profileApi "social-cloud-server/src/internal/profile/api"
	connectionApi "social-cloud-server/src/internal/connection/api"
	postApi "social-cloud-server/src/internal/post/api"

	"social-cloud-server/src/server/endpoint"
)

func (s *Server) Routes() map[string]endpoint.Handler {
	return map[string]endpoint.Handler{
		"/profile/create": profileApi.NewCreateHandler(s.Database),
		"/connection/request": connectionApi.NewRequestHandler(s.Database),
		"/post/create": postApi.NewCreateHandler(s.Database),
	}
}
