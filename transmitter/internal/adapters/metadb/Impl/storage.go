package Impl

import (
	"github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/domain"
	"github.com/jmoiron/sqlx"
)

type FilePostgres struct {
	db *sqlx.DB
}

//func (fp *FilePostgres) CreateFile(f domain.File) error {
//	//TODO implement me
//	panic("implement me")
//}

func (fp *FilePostgres) GetFileByName(name string) (*domain.File, error) {
	//TODO implement me
	panic("implement me")
}

func (fp *FilePostgres) GetFileList() ([]*domain.File, error) {
	//TODO implement me
	panic("implement me")
}

func NewFilePostgres(db *sqlx.DB) *FilePostgres {
	return &FilePostgres{db: db}
}
