package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"cloud.google.com/go/storage"

	profileModel "social-cloud-server/src/internal/profile/model"
	postModel "social-cloud-server/src/internal/post/model"
	commentModel "social-cloud-server/src/internal/comment/model"
	notificationModel "social-cloud-server/src/internal/notification/model"
	connectionModel "social-cloud-server/src/internal/connection/model"
	feedModel "social-cloud-server/src/internal/feed/model"
)

const (
	user = "postgres"
	password = "Nevergiveup1"
	dbname = "postgres"
	connection = "social-cloud-1540055012833:us-central1:social-cloud"
)

type Database struct {
	db *sql.DB
	bt *storage.BucketHandle
}

func NewDatabase() *Database {
	return &Database{}
}

func (db *Database) ConnectDatabase() error {
	psqlInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=/cloudsql/%s",
							user, password, dbname, connection)

	var err error
	db.db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	err = db.db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) BuildModels() error {
	modelQueries := []string{
		commentModel.ModelDropQuery,
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
		commentModel.ModelCreateQuery,
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
	tx, err := db.db.Begin()
	if err != nil {
		return nil, err
	}

	result, err := db.db.Exec(query)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return result, err
}

func (db *Database) ExecQuery(query string) (*sql.Rows, error) {
	tx, err := db.db.Begin()
	if err != nil {
		return nil, err
	}

	rows, err := db.db.Query(query)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return rows, err
}
