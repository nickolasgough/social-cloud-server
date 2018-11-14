package api

import (
	"context"
	"errors"
	"bytes"
	"image/png"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
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
	Username    string `json:"username"`
	Displayname string `json:"displayname"`
	Filename    string `json:"filename"`
	Imagefile   []byte `json:"imagefile"`
}

type UpdateResponse struct {
	Displayname string `json:"displayname"`
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

	imagefile, err := png.Decode(bytes.NewReader(r.Imagefile))
	if err != nil {
		return &UpdateResponse{
			Imageurl: "",
		}, err
	}

	imageurl, err := c.db.UploadImage(ctx, "pikachu.png", imagefile)
	if err != nil {
		return &UpdateResponse{
			Imageurl: "",
		}, err
	}

	_, err = c.db.ExecStatement(c.db.BuildQuery(updateQuery, r.Displayname, imageurl, r.Username))
	if err != nil {
		return UpdateResponse{
			Displayname: "",
			Imageurl: "",
		}, err
	}

	return &UpdateResponse{
		Displayname: r.Displayname,
		Imageurl: imageurl,
	}, nil
}

const updateQuery = `
UPDATE profile
SET displayname = '%s', imageurl = '%s'
WHERE username = '%s'
`
