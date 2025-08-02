package repository

import (
	"database/sql"
	"log"
)

func New() (*sql.DB, error) {
	conn, err := sql.Open("postgres", "postgres://admin:adminpassword@localhost:5432/social?sslmode=disable")
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
