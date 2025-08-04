package repository

import (
	"context"
	"database/sql"
	"log"
)

type File struct {
	ID     string
	Status string
	UserId int
}
type FilesStore struct {
	db *sql.DB
}

func (s *FilesStore) CreateFile(ctx context.Context, file *File) error {
	if err := s.db.QueryRowContext(ctx, `INSERT INTO files(user_id, status) VALUES($1, $2) returning id`, file.UserId, file.Status).Scan(&file.ID); err != nil {
		log.Fatal("Unable to create new file", err)
		return err
	}
	return nil
}

func (s *FilesStore) ChangeFileStatus(ctx context.Context, file *File) error {
	if err := s.db.QueryRowContext(ctx, `UPDATE files SET status = $2 WHERE id=$1 returning id`, file.ID, file.Status).Scan(&file.ID); err != nil {
		log.Fatal("Unable to Update file status", err)
		return err
	}
	return nil
}
