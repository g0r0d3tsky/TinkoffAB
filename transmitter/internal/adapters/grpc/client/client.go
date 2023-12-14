package client

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v9"
	"github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/domain"
	pb "github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/service"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
}
type Client struct {
	pb.TransmitterClient
}

var chunkSize = 1024 * 32

func New() error {
	cfg := &Config{}
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}
	if err := env.Parse(cfg); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	conn, err := grpc.Dial(cfg.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to establish connection: %v", err)
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
			err = GetFileList(client)
		case "file_info":
			if len(commandArgs) != 1 {
				fmt.Println("Неверное количество аргументов. Используйте: file_info <имя файла>")
				continue
			}
			err = GetFileInfo(client, commandArgs[0])
		case "get_file":
			if len(commandArgs) != 1 {
				fmt.Println("Неверное количество аргументов. Используйте: get_file <имя файла>")
				continue
			}
			err = GetFile(client, commandArgs[0])
		case "upload_file":
			if len(commandArgs) != 1 {
				fmt.Println("Неверное количество аргументов. Используйте: upload_file <путь к файлу>")
				continue
			}
			err = UploadFile(client, commandArgs[0])
		case "exit":
			fmt.Println("Выход из программы")
			return nil
		default:
			fmt.Println("Неизвестная команда")
			continue
		}

		if err != nil {
			fmt.Printf("Ошибка при выполнении команды: %v\n", err)
		}
	}
}

func UploadFile(client pb.TransmitterClient, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл: %w", err)
	}
	defer file.Close()

	stream, err := client.UploadFile(context.Background())
	if err != nil {
		return fmt.Errorf("ошибка при вызове UploadFile: %w", err)
	}

	buf := make([]byte, chunkSize)
	for {
		num, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("ошибка при чтении файла: %w", err)
		}

		if err := stream.Send(&pb.UploadFileRequest{
			Name:      filepath.Base(filePath),
			ChunkData: buf[:num],
		}); err != nil {
			return fmt.Errorf("ошибка при отправке данных файла: %w", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("ошибка при получении ответа от сервера: %w", err)
	}

	fmt.Printf("Файл успешно загружен. Имя файла на сервере: %s\n", res.FileName)

	return nil
}

func GetFileList(client pb.TransmitterClient) error {

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

func GetFileInfo(client pb.TransmitterClient, fileName string) error {

	fileInfoRequest := &pb.GetFileInfoRequest{Name: fileName}
	fileInfoResponse, err := client.GetFileInfo(context.Background(), fileInfoRequest)
	if err != nil {
		return fmt.Errorf("Ошибка при вызове GetFileInfo: %v", err)
	}

	file := fileInfoResponse.GetFile()
	fmt.Printf("Имя файла: %s\n", file.GetName())
	fmt.Printf("Размер файла: %d\n", file.GetSize())

	return nil
}

func GetFile(client pb.TransmitterClient, fileName string) error {
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
func (c *Client) Upload(file domain.FileUnit) (string, error) {

	stream, err := c.UploadFile(context.TODO())
	if err != nil {
		return "", err
	}

	// Maximum 32KB size per stream.
	buf := make([]byte, chunkSize)

	for {
		num, err := file.Payload.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		if err := stream.Send(&pb.UploadFileRequest{
			Name:      file.PayloadName,
			ChunkData: buf[:num]}); err != nil {
			return "", err
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return "", err
	}

	return res.FileName, nil
}
