package api

import (
	"context"
	"errors"
	"time"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
	"social-cloud-server/src/internal/util"
)

type AcceptHandler struct {
	db *database.Database
}

func NewAcceptHandler(db *database.Database) *AcceptHandler {
	return &AcceptHandler{
		db: db,
	}
}

type AcceptRequest struct {
	Username   string    `json:"username"`
	Connection string    `json:"connection"`
	Datetime   time.Time `json:"datetime"`
}

type AcceptResponse struct {
	Success bool `json:"success"`
}

func (c *AcceptHandler) Request() endpoint.Request {
	return &AcceptRequest{}
}

func (c *AcceptHandler) Process(ctx context.Context, request endpoint.Request) (endpoint.Response, error) {
	r, ok := request.(*AcceptRequest)
	if !ok {
		return nil, errors.New("error: received a request that is not a AcceptRequest")
	}

	lockIds := []string{"connection", "notification"}
	util.AcquireLocks(lockIds)
	defer util.ReleaseLocks(lockIds)

	_, err := c.db.ExecStatement(c.db.BuildQuery(acceptQuery, r.Username, r.Connection, r.Datetime.Format(time.RFC3339)))
	if err != nil {
		return &AcceptResponse{
			Success: false,
		}, err
	}
	_, err = c.db.ExecStatement(c.db.BuildQuery(acceptQuery, r.Connection, r.Username, r.Datetime.Format(time.RFC3339)))
	if err != nil {
		return &AcceptResponse{
			Success: false,
		}, err
	}

	_, err = c.db.ExecStatement(c.db.BuildQuery(dismissQuery, r.Username, r.Connection))
	if err != nil {
		return &AcceptResponse{
			Success: false,
		}, err
	}

	_, err = c.db.ExecStatement(c.db.BuildQuery(notifyQuery, r.Connection, r.Username, r.Datetime.Format(time.RFC3339)))
	if err != nil {
		return &AcceptResponse{
			Success: false,
		}, err
	}

	return &AcceptResponse{
		Success: true,
	}, nil
}

const acceptQuery = `
INSERT INTO connection (
	username,
	connection,
	datetime
)
VALUES (
	'%s',
	'%s',
	'%s'
);
`

const dismissQuery = `
UPDATE notification
SET dismissed = true
WHERE username = '%s' AND sender = '%s' AND type = 'connection-request';
`

const notifyQuery = `
INSERT INTO notification (
	username,
	type,
	sender,
	dismissed,
	datetime
)
VALUES (
	'%s',
	'connection-accepted',
	'%s',
	false,
	'%s'
)
`