package api

import (
	"context"
	"errors"
	"time"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
	"social-cloud-server/src/internal/notification/model"
)

type ListHandler struct {
	db *database.Database
}

func NewListHandler(db *database.Database) *ListHandler {
	return &ListHandler{
		db: db,
	}
}

type ListRequest struct {
	Username string `json:"username"`
	Cursor   string `json:"cursor"`
	Limit    string `json:"limit"`
}

type ListResponse struct {
	Data []model.Notification `json:"data"`
}

func (c *ListHandler) Request() endpoint.Request {
	return &ListRequest{}
}

func (c *ListHandler) Process(ctx context.Context, request endpoint.Request) (endpoint.Response, error) {
	cr, ok := request.(*ListRequest)
	if !ok {
		return nil, errors.New("error: received a request that is not a ListRequest")
	}

	offset := cr.Cursor
	if offset == "" {
		offset = "0"
	}
	limit := cr.Limit
	if limit == "" {
		limit = "10"
	}
	results, err := c.db.ExecQuery(c.db.BuildQuery(listQuery, cr.Username, offset, limit))
	if err != nil {
		return &ListResponse{
			Data: nil,
		}, err
	}

	var data []model.Notification
	var notification model.Notification
	var datetime string
	for results.Next() {
		err = results.Scan(&notification.Username, &notification.Type, &notification.Sender, &notification.Dismissed, &datetime)
		if err != nil {
			return &ListResponse{
				Data: nil,
			}, err
		}

		notification.Datetime, err = time.Parse(time.RFC3339, datetime)
		if err != nil {
			return ListResponse{
				Data: nil,
			}, err
		}
		data = append(data, notification)
	}

	return &ListResponse{
		Data: data,
	}, nil
}

const listQuery = `
SELECT (
	username,
	type,
	sender,
	dismissed,
	datetime
)
FROM notification
WHERE username = '%s'
OFFSET %s
LIMIT %s
`