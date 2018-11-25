package api

import (
	"context"
	"errors"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
	"social-cloud-server/src/internal/util"
	"social-cloud-server/src/bucket"
)

type UpdateHandler struct {
	db *database.Database
	b  *bucket.Bucket
}

func NewUpdateHandler(db *database.Database, b *bucket.Bucket) *UpdateHandler {
	return &UpdateHandler{
		db: db,
		b:  b,
	}
}

type UpdateRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Displayname string `json:"displayname"`
	Filename    string `json:"filename"`
	Imagefile   []byte `json:"imagefile"`
	Defaultfeed string `json:"defaultfeed"`
}

type UpdateResponse struct {
	Displayname string `json:"displayname"`
	Password    string `json:"password"`
	Imageurl    string `json:"imageurl"`
	Defaultfeed string `json:"defaultfeed"`
}

func (c *UpdateHandler) Request() endpoint.Request {
	return &UpdateRequest{}
}

func (c *UpdateHandler) Process(ctx context.Context, request endpoint.Request) (endpoint.Response, error) {
	r, ok := request.(*UpdateRequest)
	if !ok {
		return nil, errors.New("error: received a request that is not a UpdateRequest")
	}

	lockIds := []string{"profile"}
	util.AcquireLocks(lockIds)
	defer util.ReleaseLocks(lockIds)

	if r.Password != "" {
		_, err := c.db.ExecStatement(c.db.BuildQuery(passwordQuery, r.Password, r.Email))
		if err != nil {
			return &UpdateResponse{
				Displayname: "",
				Password:    "",
				Imageurl:    "",
				Defaultfeed: "",
			}, err
		}
	}

	if r.Displayname != "" {
		_, err := c.db.ExecStatement(c.db.BuildQuery(displaynameQuery, r.Displayname, r.Email))
		if err != nil {
			return &UpdateResponse{
				Displayname: "",
				Password:    "",
				Imageurl:    "",
				Defaultfeed: "",
			}, err
		}
	}

	if r.Imagefile != nil && len(r.Imagefile) > 0 {
		contentType, imagefile, err := util.DecodeImageFile(r.Filename, r.Imagefile)
		if err != nil {
			return &UpdateResponse{
				Displayname: "",
				Password:    "",
				Imageurl:    "",
				Defaultfeed: "",
			}, err
		}

		imageurl, err := c.b.UploadImage(ctx, r.Email, r.Filename, contentType, imagefile)
		if err != nil {
			return &UpdateResponse{
				Displayname: "",
				Password:    "",
				Imageurl:    "",
				Defaultfeed: "",
			}, err
		}

		_, err = c.db.ExecStatement(c.db.BuildQuery(imageurlQuery, imageurl, r.Email))
		if err != nil {
			return &UpdateResponse{
				Displayname: "",
				Password:    "",
				Imageurl:    "",
				Defaultfeed: "",
			}, err
		}
	}

	if r.Defaultfeed != "" {
		_, err := c.db.ExecStatement(c.db.BuildQuery(defaultFeedQuery, r.Defaultfeed, r.Email))
		if err != nil {
			return &UpdateResponse{
				Displayname: "",
				Password:    "",
				Imageurl:    "",
				Defaultfeed: "",
			}, err
		}
	}

	result, err := c.db.ExecQuery(c.db.BuildQuery(updateQuery, r.Email))
	if err != nil {
		return &LoginResponse{
			Displayname: "",
			Password:    "",
			Imageurl:    "",
			Defaultfeed: "",
		}, err
	}

	var ur UpdateResponse
	if result.Next() {
		err = result.Scan(&ur.Displayname, &ur.Password, &ur.Imageurl, &ur.Defaultfeed)
		if err != nil {
			return &UpdateResponse{
				Displayname: "",
				Password:    "",
				Imageurl:    "",
				Defaultfeed: "",
			}, err
		}
	}

	return &ur, nil
}

const passwordQuery = `
UPDATE profile
SET password = '%s'
WHERE email = '%s'
`

const displaynameQuery = `
UPDATE profile
SET displayname = '%s'
WHERE email = '%s'
`

const imageurlQuery = `
UPDATE profile
SET imageurl = '%s'
WHERE email = '%s'
`

const defaultFeedQuery = `
UPDATE profile
SET defaultfeed = '%s'
WHERE email = '%s'
`

const updateQuery = `
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
