package server

import (
	"context"
	miniotype "github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/adapters/minio"
	"github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/adapters/repo"
	pb "github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/service"
	"io"
)

var chunkSize = 1024 * 32

type server struct {
	miniotype.FileDB
	repo.FileMeta
	pb.UnimplementedTransmitterServer
}

//	func (serv *server) CreateFile(ctx context.Context, req *s.CreateFileRequest) (*s.SuccessResponse, error) {
//		//TODO:
//		return &s.SuccessResponse{Response: "create success"}, nil
//	}
func ServerNew(fileDB miniotype.FileDB, fileMeta repo.FileMeta) pb.TransmitterServer {
	return &server{
		FileDB:   fileDB,
		FileMeta: fileMeta,
	}
}

func (serv *server) GetFile(req *pb.GetFileRequest, stream pb.Transmitter_GetFileServer) error {
	fileName := req.GetName()

	fileUnit, err := serv.FileDB.DownloadFile(context.TODO(), fileName)
	if err != nil {
		return err
	}

	chunkBuffer := make([]byte, chunkSize)

	for {
		bytesRead, err := fileUnit.Payload.Read(chunkBuffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		response := &pb.GetFileResponse{
			Data: chunkBuffer[:bytesRead],
		}
		if err := stream.Send(response); err != nil {
			return err
		}
	}
	return nil
}

func (serv *server) GetFileList(context.Context, *pb.GetFileListRequest) (*pb.GetFileListResponse, error) {
	fileList, err := serv.FileDB.GetList()
	if err != nil {
		return nil, err
	}
	response := &pb.GetFileListResponse{
		Name: fileList,
	}

	return response, nil
}
func (serv *server) GetFileInfo(ctx context.Context, req *pb.GetFileInfoRequest) (*pb.GetFileInfoResponse, error) {
	fileName := req.GetName()
	fileInfo, err := serv.FileDB.DownloadFile(ctx, fileName)
	if err != nil {
		return nil, err
	}
	fileMeta, err := serv.FileMeta.GetFileByName(fileName)
	if err != nil {
		return nil, err
	}
	file := &pb.File{
		Name:  fileName,
		Size:  int32(fileInfo.PayloadSize),
		Owner: fileMeta.Owner,
	}
	response := &pb.GetFileInfoResponse{
		File: file,
	}

	return response, nil
}
