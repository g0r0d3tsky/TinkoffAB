package impl

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

type MinioProvider struct {
	MinioAuthData
	Client *minio.Client
}

type MinioAuthData struct {
	Url                  string `yaml:"url"`
	User                 string `yaml:"user"`
	Password             string `yaml:"password"`
	Token                string `yaml:"token"`
	Ssl                  bool   `yaml:"ssl"`
	UserObjectBucketName string `yaml:"bucket_name"`
}

//TODO: эта функция не имеет смысла, сюда написать NewMinioProvider,

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
