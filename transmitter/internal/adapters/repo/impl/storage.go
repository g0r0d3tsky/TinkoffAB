package impl

import (
	"fmt"
	"github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type FilePostgres struct {
	db *sqlx.DB
}

func (fp *FilePostgres) IsFileExists(fileName string) (bool, error) {
	query := "SELECT COUNT(*) FROM files WHERE name = $1"
	var count int
	err := fp.db.QueryRow(query, fileName).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
func (fp *FilePostgres) UploadFile(f *domain.File) (uuid.UUID, error) {
	// Проверка наличия файла в базе данных
	exists, err := fp.IsFileExists(f.Name)
	if err != nil {
		return uuid.UUID{}, err
	}
	if exists {
		return uuid.UUID{}, fmt.Errorf("file already exists in the database")
	}
	query := fmt.Sprintf("INSERT INTO files (id, name, size) values ($1, $2, $3) RETURNING id")
	row := fp.db.QueryRow(query, f.ID, f.Name, f.Size)
	if err := row.Scan(&f.ID); err != nil {
		return uuid.UUID{}, err
	}

	return f.ID, nil
}

func (fp *FilePostgres) GetFileByName(name string) (*domain.File, error) {
	file := &domain.File{}

	if err := fp.db.QueryRow(
		`SELECT id, name, size  FROM "files" f WHERE f.name = $1`, name,
	).Scan(&file.ID, &file.Name, &file.Size); err != nil {
		return nil, err
	}

	return file, nil
}

func (fp *FilePostgres) GetFileList() ([]*domain.File, error) {
	var files []*domain.File

	rows, err := fp.db.Query(
		`SELECT id, name, size FROM "files"`,
	)

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		file := &domain.File{}

		if err := rows.Scan(&file.ID, &file.Name, &file.Size); err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	return files, nil
}

func NewFilePostgres(db *sqlx.DB) *FilePostgres {
	return &FilePostgres{db: db}
}
