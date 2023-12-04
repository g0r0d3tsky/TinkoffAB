package client

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v9"
	pb "github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/service"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"strings"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
}

func New() error {

	cfg := &Config{}
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	conn, err := grpc.Dial(cfg.ServerAddress, grpc.WithTransportCredentials(
		insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Не удалось установить подключение: %v", err)
	}
	defer conn.Close()

	// Создание клиента
	client := pb.NewTransmitterClient(conn)

	for {
		fmt.Print("Введите команду: ")
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			log.Fatalf("Ошибка при чтении команды: %v", err)
		}

		args := strings.Split(input, " ")
		if len(args) == 0 {
			fmt.Println("Не указана команда")
			continue
		}

		command := args[0]
		commandArgs := args[1:]

		switch command {
		case "file_list":
			err = getFileList(client)
		case "file_info":
			err = getFileInfo(client, commandArgs)
		case "get_file":
			err = getFile(client, commandArgs)
		case "exit":
			fmt.Println("Выход из программы")
			break
		default:
			fmt.Println("Неизвестная команда")
			continue
		}

		if err != nil {
			log.Printf("Ошибка при выполнении команды: %v", err)
		}
	}
}

func getFileList(client pb.TransmitterClient) error {

	fileListRequest := &pb.GetFileListRequest{}
	fileListResponse, err := client.GetFileList(context.Background(), fileListRequest)
	if err != nil {
		return fmt.Errorf("Ошибка при вызове GetFileList: %v", err)
	}

	fmt.Println("Список файлов:")
	for _, fileName := range fileListResponse.GetName() {
		fmt.Println(fileName)
	}

	return nil
}

func getFileInfo(client pb.TransmitterClient, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("Не указано имя файла")
	}

	fileName := args[0]

	fileInfoRequest := &pb.GetFileInfoRequest{Name: fileName}
	fileInfoResponse, err := client.GetFileInfo(context.Background(), fileInfoRequest)
	if err != nil {
		return fmt.Errorf("Ошибка при вызове GetFileInfo: %v", err)
	}

	file := fileInfoResponse.GetFile()
	fmt.Printf("Имя файла: %s\n", file.GetName())
	fmt.Printf("Размер файла: %d\n", file.GetSize())
	fmt.Printf("Владелец файла: %s\n", file.GetOwner())

	return nil
}

func getFile(client pb.TransmitterClient, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("Не указано имя файла")
	}

	fileName := args[0]

	getFileRequest := &pb.GetFileRequest{Name: fileName}
	stream, err := client.GetFile(context.Background(), getFileRequest)
	if err != nil {
		return fmt.Errorf("Ошибка при вызове GetFile: %v", err)
	}

	fmt.Println("Содержимое файла:")
	for {
		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("Ошибка при получении данных файла: %v", err)
		}
		data := response.GetData()
		fmt.Print(string(data))
	}

	return nil
}
