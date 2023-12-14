package main

import (
	"github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/adapters/grpc/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"log/slog"
)

func main() {

	// Initialise gRPC connection.
	conn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}
	slog.Info("Starting dial :50051")

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("can not close connection  %v", err)
			return
		}
	}(conn)

	err = client.New()
	if err != nil {
		log.Fatalf("new client %v", err)
		return
	}
}
