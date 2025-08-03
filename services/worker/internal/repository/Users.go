package repository

import (
	"context"
	"database/sql"
)

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) CheckIfUserExists(ctx context.Context, id int64) bool {
	if err := s.db.QueryRowContext(ctx, `SELECT *  FROM USERS WHERE id= $1`, &id).Scan(); err != nil {
		return false
	}
	return true
}
