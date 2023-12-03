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
	file := &domain.File{}

	if err := fp.db.QueryRow(
		`SELECT name, size, type, owner, date  FROM "files" f WHERE f.name = $1`, name,
	).Scan(&file.Name, &file.Size, &file.Type, &file.Owner, &file.Date); err != nil {
		return nil, err
	}

	return file, nil
}

func (fp *FilePostgres) GetFileList() ([]*domain.File, error) {
	var files []*domain.File

	rows, err := fp.db.Query(
		`SELECT name, size, type, owner, date FROM "files"`,
	)

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		file := &domain.File{}

		if err := rows.Scan(&file.Name, &file.Size, &file.Type, &file.Owner, &file.Date); err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	return files, nil
}

func NewFilePostgres(db *sqlx.DB) *FilePostgres {
	return &FilePostgres{db: db}
}
