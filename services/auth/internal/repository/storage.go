package repository

import "database/sql"

type Storage struct {
	UserStore UserStore
}

func NewStorage(db *sql.DB) *Storage {

	storage := Storage{
		UserStore: UserStore{
			db: db,
		},
	}
	return &storage
}
