package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
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

func (db *Database) Connect() error {
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

func (db *Database) BuildQuery(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

func (db *Database) ExecQuery(query string) error {
	_, err := db.db.Exec(query);
	return err
}
