package repository

import (
	"database/sql"
	"log"

	"github.com/empaid/estateedge/pkg/env"
)

func New() (*sql.DB, error) {
	conn, err := sql.Open("postgres", env.GetString("POSTGRES_DB", ""))
	if err != nil {
		log.Printf("Error opening connection to database: %s", err)
	}
	return conn, nil
}

type Database struct {
	Conn *sql.DB
}

func NewConnection() (db *Database) {
	conn, err := New()
	if err != nil {
		log.Fatal("Error while connecting to database")
	}
	return &Database{
		Conn: conn,
	}
}
