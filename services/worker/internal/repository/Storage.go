package repository

import "database/sql"

type Storage struct {
	UserStore UserStore
	FileStore FilesStore
}

func NewStorage(db *sql.DB) *Storage {

	return &Storage{
		UserStore: UserStore{
			db: db,
		},
		FileStore: FilesStore{
			db: db,
		},
	}
}
