package server

import (
	"context"
	miniotype "github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/adapters/filedb"
	"github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/adapters/metadb"
	s "github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/service"
)

type server struct {
	miniotype.FileDB
	metadb.FileMeta
	s.UnimplementedTransmitterServer
}

func (serv *server) CreateFile(ctx context.Context, req *s.CreateFileRequest) (*s.SuccessResponse, error) {
	//TODO:
	return &s.SuccessResponse{Response: "create success"}, nil
}

func (serv *server) GetFile(*s.GetFileRequest, s.Transmitter_GetFileServer) error {
	// Отправьте данные клиенту через потоковую передачу (stream)
	fileData := []byte("File data") // Пример данных файла
	err := stream.Send(&s.GetFileResponse{Data: fileData})
	if err != nil {
		return err
	}

	return nil
}
func (serv *server) GetFileList(context.Context, *s.GetFileListRequest) (*s.GetFileListResponse, error) {
	// Добавьте свою логику для метода GetFileList
	// ...

	// Верните список файлов
	fileList := []*s.File{
		{FileName: "file1.txt", Size: 1024},
		{FileName: "file2.txt", Size: 2048},
	}

	return &s.GetFileListResponse{Files: fileList}, nil
}
func (serv *server) GetFileInfo(context.Context, *s.GetFileInfoRequest) (*s.GetFileInfoResponse, error) {
	// Добавьте свою логику для метода GetFileInfo
	// ...

	// Верните информацию о файле
	fileInfo := &s.FileInfo{
		FileName: "file.txt",
		Size:     1024,
		Owner:    "John Doe",
	}

	return &s.GetFileInfoResponse{Info: fileInfo}, nil
}
