package main

import (
	"context"
	"fmt"
	minio2 "github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/adapters/minio"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func main() {
	var cfg minio2.MinioAuthData
	yamlFile, _ := os.ReadFile("config_minio.yaml")
	yaml.Unmarshal(yamlFile, &cfg)
	minioClient, err := minio.New(cfg.Url, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.User, cfg.Password, ""),
		Secure: cfg.Ssl,
	})
	if err != nil {
		log.Fatal(err)
	}

	cl := &minio2.MinioProvider{
		MinioAuthData: cfg,
		Client:        minioClient,
	}
	err = cl.Connect()
	if err != nil {
		return
	}
	fmt.Println(cl.DownloadFile(context.TODO(), "file1.txt"))
}
