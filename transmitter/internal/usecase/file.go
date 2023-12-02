package usecase

import (
	"github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/domain"
	"github.com/google/uuid"
)

type FileUC interface {
	CreateFile(f domain.File) error
	GetFile(id uuid.UUID) (*domain.File, error)
	GetFileList() ([]*domain.File, error)
}
