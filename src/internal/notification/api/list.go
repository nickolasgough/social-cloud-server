package api

import (
	"context"
	"errors"
	"time"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
	"social-cloud-server/src/internal/notification/model"
	"social-cloud-server/src/internal/util"
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
	Notifications []model.Notification `json:"notifications"`
}

func (c *ListHandler) Request() endpoint.Request {
	return &ListRequest{}
}

func (c *ListHandler) Process(ctx context.Context, request endpoint.Request) (endpoint.Response, error) {
	r, ok := request.(*ListRequest)
	if !ok {
		return nil, errors.New("error: received a request that is not a ListRequest")
	}

	lockIds := []string{"notification", "profile"}
	util.AcquireLocks(lockIds)
	defer util.ReleaseLocks(lockIds)

	offset := r.Cursor
	if offset == "" {
		offset = "0"
	}
	limit := r.Limit
	if limit == "" {
		limit = "10"
	}

	results, err := c.db.ExecQuery(
		c.db.BuildQuery(
			listQuery,
			r.Username,
			offset,
			limit,
		),
	)
	if err != nil {
		return &ListResponse{
			Notifications: nil,
		}, err
	}

	var data []model.Notification
	var notification model.Notification
	var datetime string
	for results.Next() {
		err = results.Scan(
			&notification.Username,
			&notification.Type,
			&notification.Sender,
			&notification.Displayname,
			&notification.Dismissed,
			&datetime,
		)
		if err != nil {
			return &ListResponse{
				Notifications: nil,
			}, err
		}

		notification.Datetime, err = time.Parse(time.RFC3339, datetime)
		if err != nil {
			return ListResponse{
				Notifications: nil,
			}, err
		}
		data = append(data, notification)
	}

	return &ListResponse{
		Notifications: data,
	}, nil
}

const listQuery = `
SELECT
	n.username,
	n.type,
	n.sender,
	p.displayname,
	n.dismissed,
	n.datetime
FROM notification n
JOIN profile p ON p.username = n.sender
WHERE n.username = '%s' AND n.dismissed = false
ORDER BY n.datetime DESC
OFFSET %s
LIMIT %s;
`