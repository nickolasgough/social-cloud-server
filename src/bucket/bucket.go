package bucket

import (
	"context"
	"fmt"
	"image"

	"cloud.google.com/go/storage"

	"social-cloud-server/src/internal/util"
)

const (
	bucketName      = "social-cloud-1540055012833.appspot.com"
)

type Bucket struct {
	bt *storage.BucketHandle
}

func NewBucket() *Bucket {
	return &Bucket{}
}

func (b *Bucket) ConnectBucket(ctx context.Context) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	b.bt = client.Bucket(bucketName)
	return nil
}

func (b *Bucket) UploadImage(ctx context.Context, email string, filename string, contentType string, imagefile image.Image) (string, error) {
	filename = fmt.Sprintf("%s-%s", email, filename)

	object := b.bt.Object(fmt.Sprintf("%s", filename))
	writer := object.NewWriter(ctx)
	writer.ContentType = contentType

	err := util.EncodeImageFile(writer, imagefile)
	writer.Close()
	if err != nil {
		return "", err
	}

	acl := object.ACL()
	err = acl.Set(ctx, storage.AllUsers, storage.RoleReader)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, filename)
	return url, nil
}
