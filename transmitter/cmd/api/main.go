package main

import (
	server_grpc "github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/adapters/grpc/server"
	miniotype "github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/adapters/minio"
	minio2 "github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/adapters/minio/impl"
	"github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/adapters/repo"
	"github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/adapters/repo/impl"
	pb "github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/service"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
	"log"
	"log/slog"
	"net"
	"os"
)

func main() {

	var cfgMinio minio2.MinioAuthData
	yamlFile, err := os.ReadFile("config_minio.yaml")
	if err != nil {
		log.Fatalf("minio config %v", err)
		return
	}
	err = yaml.Unmarshal(yamlFile, &cfgMinio)
	if err != nil {
		log.Fatalf("can not unmarshal config %v", err)
		return
	}
	fileDB := miniotype.NewStorageMinio(cfgMinio)
	if err != nil {
		log.Fatalf("can not create minio storage %v", err)
		return
	}

	var cfgPostgres repo.ConfigPostgres
	yamlFile, _ = os.ReadFile("config_postgres.yaml")
	err = yaml.Unmarshal(yamlFile, &cfgPostgres)
	if err != nil {
		log.Fatalf("can not unmarshal postgres config")
		return
	}
	db, err := repo.NewPostgresDB(cfgPostgres)
	if err != nil {
		log.Fatalf("can not create postgres storage %v", err)
		return
	}
	metaDB := impl.NewFilePostgres(db)
	par := &server_grpc.ServerParams{
		M: fileDB,
		P: metaDB,
	}
	server := server_grpc.ServerNew(par)
	lis, err := net.Listen("tcp", ":50051")
	slog.Info("starting listen tcp :50051")
	if err != nil {
		log.Fatal(err)
	}
	defer func(lis net.Listener) {
		err := lis.Close()
		if err != nil {
			log.Fatalf("can`t cloce %w", err)
			return
		}
	}(lis)
	rpcSRV := grpc.NewServer()

	pb.RegisterTransmitterServer(rpcSRV, server)
	log.Fatal(rpcSRV.Serve(lis))
}
