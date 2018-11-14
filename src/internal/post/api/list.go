package api

import (
	"context"
	"errors"
	"time"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
	"social-cloud-server/src/internal/post/model"
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
	Feedname string `json:"feedname"`
	Cursor   string `json:"cursor"`
	Limit    string `json:"limit"`
}

type ListResponse struct {
	Posts []model.Post `json:"posts"`
}

func (c *ListHandler) Request() endpoint.Request {
	return &ListRequest{}
}

func (c *ListHandler) Process(ctx context.Context, request endpoint.Request) (endpoint.Response, error) {
	r, ok := request.(*ListRequest)
	if !ok {
		return nil, errors.New("error: received a request that is not a ListRequest")
	}

	offset := r.Cursor
	if offset == "" {
		offset = "0"
	}
	limit := r.Limit
	if limit == "" {
		limit = "25"
	}
	results, err := c.db.ExecQuery(c.db.BuildQuery(listQuery, r.Username, r.Feedname, offset, limit))
	if err != nil {
		return &ListResponse{
			Posts: nil,
		}, err
	}

	var data []model.Post
	var post model.Post
	var avator model.Avatar
	var datetime string
	for results.Next() {
		err = results.Scan(&post.Username, &avator.Displayname, &avator.Imageurl, &post.Post, &datetime)
		if err != nil {
			return &ListResponse{
				Posts: nil,
			}, err
		}
		post.Avatar = avator

		post.Datetime, err = time.Parse(time.RFC3339, datetime)
		if err != nil {
			return ListResponse{
				Posts: nil,
			}, err
		}
		data = append(data, post)
	}

	return &ListResponse{
		Posts: data,
	}, nil
}

const listQuery = `
SELECT
	po.username,
	pr.displayname,
	CASE
		WHEN pr.imageurl IS NULL THEN ''
		ELSE pr.imageurl
	END,
	po.post,
	po.datetime
FROM post po
JOIN profile pr ON pr.username = po.username
WHERE po.username IN (
	SELECT
		DISTINCT connection
	FROM feed
	WHERE username = '%s' AND feedname = '%s'
)
OFFSET %s
LIMIT %s;
`