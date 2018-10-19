package server

import (
	receiptApi "cloud-receipts/src/internal/receipt/api"
	"cloud-receipts/src/server/endpoint"
)

func (s *Server) Routes() map[string]endpoint.Handler {
	return map[string]endpoint.Handler{
		"/receipt/create": receiptApi.NewCreateHandler(),
	}
}
