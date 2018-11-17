package api

import (
	"context"
	"errors"
	"time"
	"fmt"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
	"social-cloud-server/src/internal/util"
)

type CreateHandler struct {
	db *database.Database
}

func NewCreateHandler(db *database.Database) *CreateHandler {
	return &CreateHandler{
		db: db,
	}
}

type CreateRequest struct {
	Username  string    `json:"username"`
	Post      string    `json:"post"`
	Filename  string    `json:"filename"`
	Imagefile []byte    `json:"imagefile"`
	Datetime  time.Time `json:"datetime"`
}

type CreateResponse struct {
	Success bool `json:"success"`
}

func (c *CreateHandler) Request() endpoint.Request {
	return &CreateRequest{}
}

func (c *CreateHandler) Process(ctx context.Context, request endpoint.Request) (endpoint.Response, error) {
	r, ok := request.(*CreateRequest)
	if !ok {
		return nil, errors.New("error: received a request that is not a CreateRequest")
	}

	var imageurl string
	if r.Imagefile != nil && len(r.Imagefile) > 0 {
		contentType, imagefile, err := util.DecodeImageFile(r.Filename, r.Imagefile)
		if err != nil {
			return &CreateResponse{
				Success: false,
			}, err
		}

		imageurl, err = c.db.UploadImage(ctx, r.Filename, contentType, imagefile)
		if err != nil {
			return &CreateResponse{
				Success: false,
			}, err
		}
	}
	if imageurl != "" {
		imageurl = fmt.Sprintf("'%s'", imageurl)
	} else {
		imageurl = "NULL"
	}

	_, err := c.db.ExecStatement(c.db.BuildQuery(createQuery, r.Username, r.Post, imageurl, r.Datetime.Format(time.RFC3339)))
	if err != nil {
		return &CreateResponse{
			Success: false,
		}, err
	}

	return &CreateResponse{
		Success: true,
	}, nil
}

const createQuery = `
INSERT INTO post (
	username,
	post,
	imageurl,
	likes,
	dislikes,
	datetime
)
VALUES (
	'%s',
	'%s',
	%s,
	0,
	0,
	'%s'
);
`
