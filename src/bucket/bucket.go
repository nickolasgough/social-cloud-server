package bucket

import (
	"context"
	"fmt"
	"image"

	"google.golang.org/api/option"
	"cloud.google.com/go/storage"

	"social-cloud-server/src/internal/util"
)

const (
	projectID       = "531719510691"
	bucketName      = "social-cloud-1540055012833.appspot.com"
	credentialsPath = "/home/nickolas_v_gough/projects/go/src/social-cloud-server/src/key/social-cloud-69d9b56a1450.json"
)

type Bucket struct {
	bt *storage.BucketHandle
}

func NewBucket() *Bucket {
	return &Bucket{}
}

func (b *Bucket) ConnectBucket(ctx context.Context) error {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsPath))
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
