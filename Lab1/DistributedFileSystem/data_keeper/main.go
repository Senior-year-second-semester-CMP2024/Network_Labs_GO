package main

import (
	"context"
	"log"
	"net"
	"os"

	pb "DistributedFileSystem/dfs" // Import the generated Go code

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedDFSServer
}

func (s *server) UploadFile(ctx context.Context, req *pb.UploadFileRequest) (*pb.UploadFileResponse, error) {
	// Save the uploaded file to a folder
	err := os.WriteFile(req.FileName, req.FileData, 0644)
	if err != nil {
		log.Println("Failed to save file:", err)
		return nil, err
	}
	log.Println("File saved successfully:", req.FileName)
	return &pb.UploadFileResponse{Message: "File uploaded successfully"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDFSServer(s, &server{})
	log.Println("Server started at :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
