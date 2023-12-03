package miniotype

import (
	"context"
	"github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/adapters/filedb/impl"
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

func NewStorageMinio(data impl.MinioAuthData, cl *minio.Client) *Minio {
	minioProvider := &impl.MinioProvider{
		MinioAuthData: data,
		Client:        cl,
	}

	return &Minio{
		FileDB: minioProvider,
	}
}
