package Server

import (
	"bytes"
	"context"
	miniotype "github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/adapters/minio"
	"github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/adapters/repo"
	"github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/internal/domain"
	pb "github.com/central-university-dev/2023-autumn-ab-go-hw-9-g0r0d3tsky/service"
	"github.com/google/uuid"
	"io"
	"log"
)

var chunkSize = 1024 * 32

type Server struct {
	m miniotype.FileDB
	p repo.FileMeta

	pb.UnimplementedTransmitterServer
}

type ServerParams struct {
	M miniotype.FileDB
	P repo.FileMeta
}

func ServerNew(params *ServerParams) pb.TransmitterServer {
	return &Server{
		m:                              params.M,
		p:                              params.P,
		UnimplementedTransmitterServer: pb.UnimplementedTransmitterServer{},
	}
}
func (serv *Server) UploadFile(stream pb.Transmitter_UploadFileServer) error {

	var fileData []byte
	var name string
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("no more data")
			break
		}
		if err != nil {
			return err
		}
		stream.Context()
		name = req.Name

		fileData = append(fileData, req.ChunkData...)
	}
	reader := bytes.NewReader(fileData)
	fileU := &domain.FileUnit{
		Payload:     reader,
		PayloadName: name,
		PayloadSize: len(fileData),
	}

	name, err := serv.m.UploadFile(context.TODO(), fileU)
	if err != nil {
		return err
	}
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	file := &domain.File{
		ID:   id,
		Name: fileU.PayloadName,
		Size: fileU.PayloadSize,
	}
	_, err = serv.p.UploadFile(file)
	if err != nil {
		return err
	}

	return nil

}

func (serv *Server) GetFile(req *pb.GetFileRequest, stream pb.Transmitter_GetFileServer) error {
	fileName := req.GetName()

	fileUnit, err := serv.m.DownloadFile(context.TODO(), fileName)
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

func (serv *Server) GetFileList(context.Context, *pb.GetFileListRequest) (*pb.GetFileListResponse, error) {
	fileList, err := serv.m.GetList()
	if err != nil {
		return nil, err
	}
	response := &pb.GetFileListResponse{
		Name: fileList,
	}

	return response, nil
}
func (serv *Server) GetFileInfo(ctx context.Context, req *pb.GetFileInfoRequest) (*pb.GetFileInfoResponse, error) {
	fileName := req.GetName()
	fileInfo, err := serv.m.DownloadFile(ctx, fileName)
	if err != nil {
		return nil, err
	}
	file := &pb.File{
		Name: fileName,
		Size: int32(fileInfo.PayloadSize),
	}
	response := &pb.GetFileInfoResponse{
		File: file,
	}

	return response, nil
}
