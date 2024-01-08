package main

import (
	server_grpc "github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/adapters/grpc/server"
	miniotype "github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/adapters/minio"
	"github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/adapters/repo"
	"github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/adapters/repo/impl"
	"github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/config"
	pb "github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/service"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"log/slog"
	"net"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("config read: %v", err)
		return
	}

	fileDB := miniotype.NewStorageMinio(cfg.Minio)
	if err != nil {
		log.Fatalf("can not create minio storage %v", err)
		return
	}

	db, err := repo.NewPostgresDB(cfg.Postgres)
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
			log.Fatalf("can`t close %w", err)
			return
		}
	}(lis)
	rpcSRV := grpc.NewServer()

	pb.RegisterTransmitterServer(rpcSRV, server)
	log.Fatal(rpcSRV.Serve(lis))
}
