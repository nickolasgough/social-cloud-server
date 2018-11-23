package api

import (
	"context"
	"errors"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
	"social-cloud-server/src/internal/util"
)

type UpdateHandler struct {
	db *database.Database
}

func NewUpdateHandler(db *database.Database) *UpdateHandler {
	return &UpdateHandler{
		db: db,
	}
}

type UpdateRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Displayname string `json:"displayname"`
	Filename    string `json:"filename"`
	Imagefile   []byte `json:"imagefile"`
}

type UpdateResponse struct {
	Displayname string `json:"displayname"`
	Password    string `json:"password"`
	Imageurl    string `json:"imageurl"`
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
			}, err
		}

		imageurl, err := c.db.UploadImage(ctx, r.Email, r.Filename, contentType, imagefile)
		if err != nil {
			return &UpdateResponse{
				Displayname: "",
				Password:    "",
				Imageurl:    "",
			}, err
		}

		_, err = c.db.ExecStatement(c.db.BuildQuery(imageurlQuery, imageurl, r.Email))
		if err != nil {
			return &UpdateResponse{
				Displayname: "",
				Password:    "",
				Imageurl:    "",
			}, err
		}
	}

	result, err := c.db.ExecQuery(c.db.BuildQuery(updateQuery, r.Email))
	if err != nil {
		return &LoginResponse{
			Displayname: "",
			Password:    "",
			Imageurl:    "",
		}, err
	}

	var ur UpdateResponse
	if result.Next() {
		err = result.Scan(&ur.Displayname, &ur.Password, &ur.Imageurl)
		if err != nil {
			return &UpdateResponse{
				Displayname: "",
				Password:    "",
				Imageurl:    "",
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

const updateQuery = `
SELECT
	displayname,
	password,
	CASE 
		WHEN imageurl IS NULL THEN ''
		ELSE imageurl
	END
FROM profile
WHERE email = '%s';
`