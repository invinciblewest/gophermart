package postgres

import (
	"database/sql"
)

type PGRepository struct {
	db *sql.DB
}

func NewPGRepository(db *sql.DB) *PGRepository {
	return &PGRepository{
		db: db,
	}
}
