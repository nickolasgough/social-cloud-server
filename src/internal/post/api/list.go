package api

import (
	"context"
	"errors"
	"time"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
	"social-cloud-server/src/internal/post/model"
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
	Email    string `json:"email"`
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

	lockIds := []string{"post", "profile", "feed"}
	util.AcquireLocks(lockIds)
	defer util.ReleaseLocks(lockIds)

	offset := r.Cursor
	if offset == "" {
		offset = "0"
	}
	limit := r.Limit
	if limit == "" {
		limit = "25"
	}

	results, err := c.db.ExecQuery(
		c.db.BuildQuery(
			listQuery,
			r.Email,
			r.Email,
			r.Email,
			r.Feedname,
			offset,
			limit,
		),
	)
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
		err = results.Scan(
			&post.Email,
			&avator.Displayname,
			&avator.Imageurl,
			&post.Post,
			&post.Imageurl,
			&post.Linkurl,
			&post.Likes,
			&post.Dislikes,
			&post.Liked,
			&post.Disliked,
			&datetime,
		)
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
	po.email,
	pr.displayname,
	CASE
		WHEN pr.imageurl IS NULL THEN ''
		ELSE pr.imageurl
	END,
	po.post,
	CASE
		WHEN po.imageurl IS NULL THEN ''
		ELSE po.imageurl
	END,
	po.linkurl,
	po.likes,
	po.dislikes,
	CASE
		WHEN (
			SELECT
				COUNT(r.connection)
			FROM reaction r
			WHERE r.email = po.email AND r.posttime = po.datetime AND r.connection = '%s' AND r.reaction = 'liked'
		) > 0 THEN TRUE
		ELSE FALSE
	END,
	CASE
		WHEN (
			SELECT
				COUNT(r.connection)
			FROM reaction r
			WHERE r.email = po.email AND r.posttime = po.datetime AND r.connection = '%s' AND r.reaction = 'disliked'
		) > 0 THEN TRUE
		ELSE FALSE
	END,
	po.datetime
FROM post po
JOIN profile pr ON pr.email = po.email
WHERE po.email IN (
	SELECT
		DISTINCT connection
	FROM feed f
	WHERE f.email = '%s' AND f.feedname = '%s'
)
ORDER BY po.datetime DESC
OFFSET %s
LIMIT %s;
`
