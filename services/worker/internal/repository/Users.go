package repository

import (
	"context"
	"database/sql"
	"log"
)

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) CheckIfUserExists(ctx context.Context, id int64) bool {
	userExists := false
	if err := s.db.QueryRowContext(ctx, `SELECT EXISTS( SELECT 1 FROM USERS WHERE id=$1) `, &id).Scan(&userExists); err != nil {
		log.Print("error while checking for user: ", err)
		return false
	}
	return userExists
}
