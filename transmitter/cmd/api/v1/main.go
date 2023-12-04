package main

import (
	"github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/adapters/grpc/client"
	server_grpc "github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/adapters/grpc/server"
	miniotype "github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/adapters/minio"
	minio2 "github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/adapters/minio/impl"
	"github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/adapters/repo"
	"github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/adapters/repo/impl"
	"gopkg.in/yaml.v3"
	"os"
)

func main() {
	var cfgMinio minio2.MinioAuthData
	yamlFile, err := os.ReadFile("config_minio.yaml")
	if err != nil {
		return
	}
	err = yaml.Unmarshal(yamlFile, &cfgMinio)
	if err != nil {
		return
	}
	fileDB := miniotype.NewStorageMinio(cfgMinio)
	if err != nil {
		return
	}

	var cfgPostgres repo.ConfigPostgres
	yamlFile, _ = os.ReadFile("config_postgres.yaml")
	err = yaml.Unmarshal(yamlFile, &cfgPostgres)
	if err != nil {
		return
	}
	db, err := repo.NewPostgresDB(cfgPostgres)
	if err != nil {
		return
	}
	metaDB := impl.NewFilePostgres(db)
	server := server_grpc.ServerNew(fileDB, metaDB)
	err := client.New()
	if err != nil {
		// Обработка ошибки
	}
}
