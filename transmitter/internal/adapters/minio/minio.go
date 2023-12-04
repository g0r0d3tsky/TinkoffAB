package miniotype

import (
	"context"
	"github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/adapters/minio/impl"
	"github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/domain"
	"github.com/minio/minio-go/v7"
)

type FileDB interface {
	DownloadFile(ctx context.Context, name string) (*domain.FileUnit, error)
	GetList() ([]string, error)
}

type Minio struct {
	FileDB
}

func NewStorageMinio(data impl.MinioAuthData) *Minio {
	minioProvider := &impl.MinioProvider{
		MinioAuthData: data,
		Client:        &minio.Client{},
	}
	err := minioProvider.Connect()
	if err != nil {
		return nil
	}
	return &Minio{
		FileDB: minioProvider,
	}
}
