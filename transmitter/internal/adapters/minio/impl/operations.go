package impl

import (
	"bytes"
	"context"
	"github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/domain"
	"github.com/minio/minio-go/v7"
	"io"
	"log"
)

// UploadFile - Отправляет файл в minio
//func (m *MinioProvider) UploadFile(ctx context.Context, object domain.File) (string, error) {
//
//	_, err := m.client.PutObject(
//		ctx,
//		UserObjectsBucketName, // Константа с именем бакета
//		object.Name,
//		object.ID,
//		object.P,
//		minio.PutObjectOptions{ContentType: "image/png"},
//	)
//
//	return imageName, err
//}

// DownloadFile - Возвращает файл из minio
func (m *MinioProvider) DownloadFile(ctx context.Context, name string) (*domain.FileUnit, error) {
	reader, err := m.Client.GetObject(
		ctx,
		m.UserObjectBucketName,
		name,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	fileUnit := &domain.FileUnit{
		Payload:     bytes.NewReader(data),
		PayloadName: name,           // Имя файла (можно использовать значение переменной name)
		PayloadSize: int(len(data)), // Размер файла (используем длину данных)
	}

	return fileUnit, nil
}
func (m *MinioProvider) GetList() ([]string, error) {
	// List objects in the bucket
	objectCh := m.Client.ListObjects(context.Background(), m.UserObjectBucketName, minio.ListObjectsOptions{
		Recursive: true, // Set to false if you only want to list objects in the root directory
	})

	var list []string
	for object := range objectCh {
		if object.Err != nil {
			log.Println(object.Err)
			continue
		}
		list = append(list, object.Key)
	}
	return list, nil
}
