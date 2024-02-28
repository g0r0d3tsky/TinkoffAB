package main

import (
	"2024-spring-ab-go-hw-1-template-g0r0d3tsky/chat/internal/config"
	"2024-spring-ab-go-hw-1-template-g0r0d3tsky/chat/internal/repository"
	"2024-spring-ab-go-hw-1-template-g0r0d3tsky/chat/internal/usecase"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
)

type app struct {
	UC     *usecase.UC
	config *config.Config
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	c, err := config.Read()

	if err != nil {
		log.Println("failed to read config:", err.Error())
		return
	}

	dbPool, err := repository.Connect(c)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer func() {
		if dbPool != nil {
			dbPool.Close()
		}
	}()
	repo := repository.New(dbPool)
	//fmt.Printf("%+v", repo)
	service := usecase.New(repo)
	//fmt.Printf("%+v", service)

	app := &app{
		UC:     service,
		config: c,
	}
	err = app.Serve()
	if err != nil {
		slog.Error("running server", err)
	}
}