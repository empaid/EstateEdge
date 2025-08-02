package repository

import (
	"context"
	"database/sql"
	"log"
)

type User struct {
	Username string
	ID       string
	Password []byte
	Email    string
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) GetUser(ctx context.Context, username string) (*User, error) {
	user := User{}
	err := s.db.QueryRowContext(ctx, `SELECT id, email, username, password FROM USERS WHERE username=$1`, username).Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	return &user, err
}

func (s *UserStore) RegisterUser(ctx context.Context, user *User) error {
	err := s.db.QueryRowContext(ctx, `INSERT INTO USERS(username, email, password) VALUES($1, $2, $3) returning id;`, user.Username, user.Email, user.Password).Scan(
		&user.ID,
	)
	if err != nil {
		log.Fatal("Error while adding user", err)
	}
	return err
}
