package api

import (
	"context"
	"errors"
	"time"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
	"social-cloud-server/src/internal/util"
)

type GoogleHandler struct {
	db *database.Database
}

func NewGoogleHandler(db *database.Database) *GoogleHandler {
	return &GoogleHandler{
		db: db,
	}
}

type GoogleRequest struct {
	Email       string    `json:"email"`
	DisplayName string    `json:"displayname"`
	Imageurl    string    `json:"imageurl"`
	Datetime    time.Time `json:"datetime"`
}

type GoogleResponse struct {
	Displayname string `json:"displayname"`
	Password    string `json:"password"`
	Imageurl    string `json:"imageurl"`
	Defaultfeed string `json:"defaultfeed"`
}

func (c *GoogleHandler) Request() endpoint.Request {
	return &GoogleRequest{}
}

func (c *GoogleHandler) Process(ctx context.Context, request endpoint.Request) (endpoint.Response, error) {
	r, ok := request.(*GoogleRequest)
	if !ok {
		return nil, errors.New("error: received a request that is not a GoogleRequest")
	}

	lockIds := []string{"profile"}
	util.AcquireLocks(lockIds)
	defer util.ReleaseLocks(lockIds)

	c.db.ExecQuery(c.db.BuildQuery(googleCreateQuery, r.Email, r.DisplayName, r.Imageurl, r.Datetime.Format(time.RFC3339)))

	result, err := c.db.ExecQuery(c.db.BuildQuery(googleSignInQuery, r.Email))
	if err != nil {
		return &GoogleResponse{
			Displayname: "",
			Password:    "",
			Imageurl:    "",
			Defaultfeed: "",
		}, err
	}

	var gr GoogleResponse
	if result.Next() {
		err = result.Scan(&gr.Displayname, &gr.Password, &gr.Imageurl, &gr.Defaultfeed)
		if err != nil {
			return &GoogleResponse{
				Displayname: "",
				Password:    "",
				Imageurl:    "",
				Defaultfeed: "",
			}, err
		}
	}

	return &gr, nil
}

const googleCreateQuery = `
INSERT INTO profile (
	email,
	password,
	displayname,
	imageurl,
	defaultfeed,
	datetime
)
VALUES (
	'%s',
	'default-password',
	'%s',
	'%s',
	NULL,
	'%s'
);
`

const googleSignInQuery = `
SELECT
	displayname,
	password,
	CASE 
		WHEN imageurl IS NULL THEN ''
		ELSE imageurl
	END,
	CASE 
		WHEN defaultfeed IS NULL THEN ''
		ELSE defaultfeed
	END
FROM profile
WHERE email = '%s';
`
