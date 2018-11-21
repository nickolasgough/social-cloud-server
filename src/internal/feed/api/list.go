package api

import (
	"context"
	"errors"
	"time"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
	"social-cloud-server/src/internal/feed/model"
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
	Feeds []model.Feed `json:"feeds"`
}

func (c *ListHandler) Request() endpoint.Request {
	return &ListRequest{}
}

func (c *ListHandler) Process(ctx context.Context, request endpoint.Request) (endpoint.Response, error) {
	r, ok := request.(*ListRequest)
	if !ok {
		return nil, errors.New("error: received a request that is not a ListRequest")
	}

	lockIds := []string{"feed"}
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
	results, err := c.db.ExecQuery(c.db.BuildQuery(listQuery, r.Username, offset, limit))
	if err != nil {
		return &ListResponse{
			Feeds: nil,
		}, err
	}

	fmap := make(map[string]model.Feed)
	var temp model.Feed
	var feed model.Feed
	var member model.Member
	var fdatetime string
	var mdatetime string;
	for results.Next() {
		err = results.Scan(&feed.Username, &feed.Feedname, &member.Connection, &mdatetime, &fdatetime)
		if err != nil {
			return &ListResponse{
				Feeds: nil,
			}, err
		}

		member.Datetime, err = time.Parse(time.RFC3339, mdatetime)
		if err != nil {
			return ListResponse{
				Feeds: nil,
			}, err
		}
		feed.Datetime, err = time.Parse(time.RFC3339, fdatetime)
		if err != nil {
			return ListResponse{
				Feeds: nil,
			}, err
		}

		temp = fmap[feed.Feedname]
		temp.Username = feed.Username
		temp.Feedname = feed.Feedname
		temp.Datetime = feed.Datetime
		temp.Members = append(temp.Members, member)
		fmap[temp.Feedname] = temp
	}

	var feeds []model.Feed
	for _, feed := range fmap {
		feeds = append(feeds, feed)
	}

	return &ListResponse{
		Feeds: feeds,
	}, nil
}

const listQuery = `
SELECT
	username,
	feedname,
	connection,
	joined,
	datetime
FROM feed
WHERE username = '%s'
ORDER BY feedname ASC
OFFSET %s
LIMIT %s;
`
