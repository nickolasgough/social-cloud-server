package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	connectionModel "social-cloud-server/src/internal/connection/model"
	notificationModel "social-cloud-server/src/internal/notification/model"
	postModel "social-cloud-server/src/internal/post/model"
	profileModel "social-cloud-server/src/internal/profile/model"
)

const (
	host = "localhost"
	port = 5432
	user = "postgres"
	password = "Nevergiveup1"
	dbname = "socialclouddb"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() *Database {
	return &Database{}
}

func (db *Database) ConnectDatabase() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
							host, port, user, password, dbname)

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

func (db *Database) BuildModels() error {
	modelQueries := []string{
		postModel.ModelDropQuery,
		connectionModel.ModelDropQuery,
		notificationModel.ModelDropQuery,
		profileModel.ModelDropQuery,

		profileModel.ModelCreateQuery,
		notificationModel.ModelCreateQuery,
		connectionModel.ModelCreateQuery,
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
