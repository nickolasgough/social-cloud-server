package server

import (
	"social-cloud-server/src/server/endpoint"

	profileApi "social-cloud-server/src/internal/profile/api"
	connectionApi "social-cloud-server/src/internal/connection/api"
	postApi "social-cloud-server/src/internal/post/api"
	commentApi "social-cloud-server/src/internal/comment/api"
	notificationApi "social-cloud-server/src/internal/notification/api"
	feedApi "social-cloud-server/src/internal/feed/api"
)

func (s *Server) Routes() map[string]endpoint.Handler {
	return map[string]endpoint.Handler{
		"/profile/create": profileApi.NewCreateHandler(s.Database),
		"/profile/login": profileApi.NewLoginHandler(s.Database),
		"/profile/update": profileApi.NewUpdateHandler(s.Database, s.Bucket),
		"/profile/google": profileApi.NewGoogleHandler(s.Database),
		"/user/search": profileApi.NewSearchHandler(s.Database),
		"/connection/request": connectionApi.NewRequestHandler(s.Database),
		"/connection/accept": connectionApi.NewAcceptHandler(s.Database),
		"/connection/decline": connectionApi.NewDeclineHandler(s.Database),
		"/connection/list": connectionApi.NewListHandler(s.Database),
		"/post/create": postApi.NewCreateHandler(s.Database, s.Bucket),
		"/post/react": postApi.NewReactHandler(s.Database),
		"/post/list": postApi.NewListHandler(s.Database),
		"/comment/create": commentApi.NewCreateHandler(s.Database),
		"/comment/list": commentApi.NewListHandler(s.Database),
		"/notification/list": notificationApi.NewListHandler(s.Database),
		"/notification/dismiss": notificationApi.NewDismissHandler(s.Database),
		"/feed/create": feedApi.NewCreateHandler(s.Database),
		"/feed/list": feedApi.NewListHandler(s.Database),
	}
}
