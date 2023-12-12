package repo

import (
	"github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/adapters/repo/impl"
	"github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type FileMeta interface {
	UploadFile(f *domain.File) (uuid.UUID, error)
	GetFileByName(name string) (*domain.File, error)
	GetFileList() ([]*domain.File, error)
}
type Storage struct {
	FileMeta
}

func NewStoragePostgres(db *sqlx.DB) *Storage {
	return &Storage{
		FileMeta: impl.NewFilePostgres(db),
	}
}
