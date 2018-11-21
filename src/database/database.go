package database

import (
	"database/sql"
	"fmt"
	"context"
	"image"
	"google.golang.org/api/option"

	_ "github.com/lib/pq"
	"cloud.google.com/go/storage"

	profileModel "social-cloud-server/src/internal/profile/model"
	postModel "social-cloud-server/src/internal/post/model"
	notificationModel "social-cloud-server/src/internal/notification/model"
	connectionModel "social-cloud-server/src/internal/connection/model"
	feedModel "social-cloud-server/src/internal/feed/model"

	"social-cloud-server/src/internal/util"
)

const (
	host = "35.202.106.171"
	port = 5432
	user = "postgres"
	password = "Nevergiveup1"
	dbname = "postgres"

	projectID = "531719510691"
	bucketName = "social-cloud-1540055012833.appspot.com"
	credentialsPath = "/home/nickolas_v_gough/projects/go/src/social-cloud-server/src/key/social-cloud-69d9b56a1450.json"
)

type Database struct {
	db *sql.DB
	bt *storage.BucketHandle
}

func NewDatabase() *Database {
	return &Database{}
}

func (db *Database) ConnectDatabase() error {
	psqlInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
							host, user, password, dbname)

	var err error
	db.db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	err = db.db.Ping()
	if err != nil {
		return err
	}
	fmt.Printf("Successfully established a connection to the database %s\n", dbname)

	return nil
}

func (db *Database) ConnectBucket(ctx context.Context) error {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		return err
	}

	db.bt = client.Bucket(bucketName)
	return nil
}

func (db *Database) BuildModels() error {
	modelQueries := []string{
		//postModel.ModelDropQuery,
		//feedModel.ModelDropQuery,
		//connectionModel.ModelDropQuery,
		//notificationModel.ModelDropQuery,
		//profileModel.ModelDropQuery,

		profileModel.ModelCreateQuery,
		notificationModel.ModelCreateQuery,
		connectionModel.ModelCreateQuery,
		feedModel.ModelCreateQuery,
		postModel.ModelCreateQuery,
	}

	for _, modelQuery := range modelQueries {
		_, err := db.ExecStatement(modelQuery)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *Database) BuildQuery(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

func (db *Database) ExecStatement(query string) (sql.Result, error) {
	return db.db.Exec(query)
}

func (db *Database) ExecQuery(query string) (*sql.Rows, error) {
	return db.db.Query(query)
}

func (db *Database) UploadImage(ctx context.Context, username string, filename string, contentType string, imagefile image.Image) (string, error) {
	object := db.bt.Object(fmt.Sprintf("%s-%s", username, filename))
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
