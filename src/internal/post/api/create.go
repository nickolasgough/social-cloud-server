package api

import (
	"context"
	"errors"
	"time"
	"fmt"

	"social-cloud-server/src/server/endpoint"
	"social-cloud-server/src/database"
	urlShortener "social-cloud-server/src/url-shortener"
	"social-cloud-server/src/internal/util"
	"social-cloud-server/src/bucket"
)

type CreateHandler struct {
	db *database.Database
	b  *bucket.Bucket
}

func NewCreateHandler(db *database.Database, b *bucket.Bucket) *CreateHandler {
	return &CreateHandler{
		db: db,
		b:  b,
	}
}

type CreateRequest struct {
	Email     string    `json:"email"`
	Post      string    `json:"post"`
	Filename  string    `json:"filename"`
	Imagefile []byte    `json:"imagefile"`
	Linkurl   string    `json:"linkurl"`
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

	lockIds := []string{"post"}
	util.AcquireLocks(lockIds)
	defer util.ReleaseLocks(lockIds)

	var imageurl string
	if r.Imagefile != nil && len(r.Imagefile) > 0 {
		contentType, imagefile, err := util.DecodeImageFile(r.Filename, r.Imagefile)
		if err != nil {
			return &CreateResponse{
				Success: false,
			}, err
		}

		imageurl, err = c.b.UploadImage(ctx, r.Email, r.Filename, contentType, imagefile)
		if err != nil {
			return &CreateResponse{
				Success: false,
			}, err
		}

		imageurl = fmt.Sprintf("'%s'", imageurl)
	} else {
		imageurl = "NULL"
	}

	var linkurl string
	if r.Linkurl != "" {
		var err error
		linkurl, err = urlShortener.ShortenUrl(r.Linkurl)
		if err != nil {
			return &CreateResponse{
				Success: false,
			}, err
		}

		linkurl = fmt.Sprintf("'%s'", linkurl)
	} else {
		linkurl = "NULL"
	}

	_, err := c.db.ExecStatement(
		c.db.BuildQuery(
			createQuery,
			r.Email,
			r.Post,
			imageurl,
			linkurl,
			r.Datetime.Format(time.RFC3339),
		),
	)
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
	email,
	post,
	imageurl,
	linkurl,
	likes,
	dislikes,
	datetime
)
VALUES (
	'%s',
	'%s',
	%s,
	%s,
	0,
	0,
	'%s'
);
`
