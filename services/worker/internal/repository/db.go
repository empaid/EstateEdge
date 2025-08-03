package repository

import (
	"database/sql"
	"log"

	"github.com/empaid/estateedge/pkg/env"
	_ "github.com/lib/pq"
)

func NewDBConection() (*sql.DB, error) {
	db, err := sql.Open("postgres", env.GetString("POSTGRES_DB", ""))
	if err != nil {
		log.Print("Not able to connect to the database ", err)
		return nil, err
	}
	return db, nil
}
