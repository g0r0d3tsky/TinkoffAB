package repository

import (
	"github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/domain"
	"github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/repository/Impl"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type File interface {
	CreateFile(f domain.File) error
	GetFile(id uuid.UUID) (*domain.File, error)
	GetFileList() ([]*domain.File, error)
}
type Storage struct {
	File
}

func NewStoragePostgres(db *sqlx.DB) *Storage {
	return &Storage{
		File: Impl.NewFilePostgres(db),
	}
}
