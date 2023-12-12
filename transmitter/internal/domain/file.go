package domain

import (
	"github.com/google/uuid"
	"io"
)

type File struct {
	ID   uuid.UUID
	Name string
	Size int
}

type FileUnit struct {
	Payload     io.Reader
	PayloadName string
	PayloadSize int
}
