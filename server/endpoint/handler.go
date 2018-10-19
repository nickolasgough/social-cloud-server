package endpoint

import "context"

type Handler interface {
	Request() Request
	Process(c context.Context, r Request) (Response, error)
}

type Request interface {}

type Response interface {}
