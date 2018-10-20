package api

import (
	"context"

	"social-cloud-server/src/server/endpoint"
)

type CreateHandler struct {}

func NewCreateHandler() *CreateHandler {
	return &CreateHandler{}
}

type CreateRequest struct {
	Name string `json:"name"`
}

type CreateResponse struct {
	Success bool
}

func (c *CreateHandler) Request() endpoint.Request {
	return &CreateRequest{}
}

func (c *CreateHandler) Process(ctx context.Context, request endpoint.Request) (endpoint.Response, error) {
	return &CreateResponse{
		Success: true,
	}, nil
}