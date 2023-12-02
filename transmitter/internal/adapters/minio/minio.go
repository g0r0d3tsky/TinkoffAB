package minio

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

// MinioProvider - Наш провайдер для хранилища
type MinioProvider struct {
	MinioAuthData
	Client *minio.Client
}

type MinioAuthData struct {
	Url      string `yaml:"url"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	//Token                string `yaml:"token"`
	Ssl                  bool   `yaml:"ssl"`
	UserObjectBucketName string `yaml:"bucket_name"`
}

func (m *MinioProvider) Connect() error {
	var err error
	m.Client, err = minio.New(m.Url, &minio.Options{
		Creds:  credentials.NewStaticV4(m.User, m.Password, ""),
		Secure: m.Ssl,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return err
}
