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
		postModel.ModelDropQuery,
		feedModel.ModelDropQuery,
		connectionModel.ModelDropQuery,
		notificationModel.ModelDropQuery,
		profileModel.ModelDropQuery,

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
	db.begin()
	result, err := db.db.Exec(query)
	if err != nil {
		db.rollback()
	} else {
		db.commit()
	}
	return result, err
}

func (db *Database) ExecQuery(query string) (*sql.Rows, error) {
	db.begin()
	rows, err := db.db.Query(query)
	if err != nil {
		db.rollback()
	} else {
		db.commit()
	}
	return rows, err
}

func (db *Database) begin() {
	if _, err := db.db.Exec("BEGIN;"); err != nil {
		fmt.Printf("Begin errored with: %s\n", err.Error())
	}
	fmt.Println("Beginning transaction")
}

func (db *Database) commit() {
	fmt.Println("Committing transaction")
	if _, err := db.db.Exec("COMMIT;"); err != nil {
		fmt.Printf("Commit errored with: %s\n", err.Error())
	}
}

func (db *Database) rollback() {
	fmt.Println("Restarting transaction")
	if _, err := db.db.Exec("ROLLBACK;"); err != nil {
		fmt.Printf("Rollback errored with: %s\n", err.Error())
	}
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
