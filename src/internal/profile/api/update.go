package api

import (
	"context"
	"errors"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
	"social-cloud-server/src/internal/util"
	"fmt"
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

	var imageurl string
	if r.Imagefile != nil && len(r.Imagefile) > 0 {
		contentType, imagefile, err := util.DecodeImageFile(r.Filename, r.Imagefile)
		if err != nil {
			return &UpdateResponse{
				Displayname: "",
				Password:    "",
				Imageurl:    "",
			}, err
		}

		imageurl, err = c.db.UploadImage(ctx, r.Email, r.Filename, contentType, imagefile)
		if err != nil {
			return &UpdateResponse{
				Displayname: "",
				Password:    "",
				Imageurl:    "",
			}, err
		}
	}

	var newurl string
	if imageurl != "" {
		newurl = fmt.Sprintf("'%s'", imageurl)
	} else {
		newurl = "NULL"
	}

	_, err := c.db.ExecStatement(c.db.BuildQuery(updateQuery, r.Displayname, r.Password, newurl, r.Email))
	if err != nil {
		return UpdateResponse{
			Displayname: "",
			Password:    "",
			Imageurl:    "",
		}, err
	}

	if imageurl == "NULL" {
		imageurl = ""
	}

	return &UpdateResponse{
		Displayname: r.Displayname,
		Password:    r.Password,
		Imageurl:    imageurl,
	}, nil
}

const updateQuery = `
UPDATE profile
SET displayname = '%s', password = '%s', imageurl = %s
WHERE email = '%s'
`
