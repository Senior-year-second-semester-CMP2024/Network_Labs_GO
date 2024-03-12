package main

import (
	"context"
	"log"
	"net"
	"os"
	"sync"

	pb "DistributedFileSystem/dfs" // Import the generated Go code

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedDFSServer
	client pb.DFSClient
}

func (s *server) UploadFile(ctx context.Context, req *pb.UploadFileRequest) (*pb.Empty, error) {
	// Save the uploaded file to a folder
	err := os.WriteFile("./data_keeper/"+req.FileName, req.FileData, 0644)
	if err != nil {
		log.Println("Failed to save file:", err)
		return nil, err
	}
	log.Println("File saved successfully:", req.FileName)
	// Call the UploadSuccess RPC
	err = s.callUploadSuccess(req.FileName, "DataKeeperNodeName", "FilePathOnNode")
	if err != nil {
		log.Println("Failed to call UploadSuccess:", err)
		// Handle error if necessary
	}
	return &pb.Empty{}, nil
}
func (s *server) DownloadFile(ctx context.Context, req *pb.DownloadFileRequest) (*pb.DownloadFileResponse, error) {
	// Load the file
	file, err := os.ReadFile(req.FileName)
	if err != nil {
		log.Println("Failed to load file:", err)
		return nil, err
	}
	res := &pb.DownloadFileResponse{
		FileData: file,
	}
	return res, nil
}
func (s *server) callUploadSuccess(fileName string, nodeName string, filePath string) error {
	// Prepare the request
	request := &pb.UploadSuccessRequest{
		FileName:           fileName,
		DataKeeperNodeName: nodeName,
		FilePathOnNode:     filePath,
	}

	// Call the UploadSuccess RPC using the client
	_, err := s.client.UploadSuccess(context.Background(), request)
	if err != nil {
		return err
	}

	return nil
}
func PingMasterTracker(client pb.DFSClient) error {
	// Prepare the request
	// TODO: make the name variable
	// TODO: add available ports to the request
	req := &pb.PingMasterTrackerRequest{
		NodeName: "node1",
	}
	// Call the PingMasterTracker RPC on the master tracker node
	_, err := client.PingMasterTracker(context.Background(), req)
	if err != nil {
		log.Println("Failed to ping master tracker node:", err)
		return err
	}

	return nil
}
func main() {
	// Client setup
	// Set up a gRPC connection to the server implementing UploadSuccess
	ClientConn, err := grpc.Dial("localhost:8080", grpc.WithInsecure()) // Update with actual server address
	if err != nil {
		log.Fatalf("failed to connect to data keeper: %v", err)
	}
	defer ClientConn.Close()

	// Create a client for the UploadSuccess service
	client := pb.NewDFSClient(ClientConn)

	// Server setup
	ports := []string{":50051", ":50052", ":50053"}
	var wg sync.WaitGroup
	wg.Add(len(ports))

	for _, port := range ports {
		go startServer(port, &wg, client)
	}

	wg.Wait()
}
func startServer(port string, wg *sync.WaitGroup, client pb.DFSClient) {
	defer wg.Done()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	s := grpc.NewServer()
	srv := &server{
		client: client,
	}
	pb.RegisterDFSServer(s, srv)
	log.Printf("Server started at %s\n", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
