package main

import (
	"context"
	"fmt"
	"net"

	pb "DistributedFileSystem/dfs" // Import generated proto file

	"google.golang.org/grpc"
)

type FileRecord struct {
	FileName        string
	DataKeeperNode  string
	FilePath        string
	IsDataNodeAlive bool
}

// MasterTrackerServer implements the DFS service
type MasterTrackerServer struct {
	pb.UnimplementedDFSServer
	lookupTable []FileRecord
}

func (s *MasterTrackerServer) RequestToUpload(ctx context.Context, req *pb.Empty) (*pb.RequestToUploadResponse, error) {
	// Implement logic to handle client request to upload a file
	// token = port number of the data keeper node that exist in the lookup table
	token := "50051"
	return &pb.RequestToUploadResponse{
		Token: token,
	}, nil
}

func (s *MasterTrackerServer) NotifyClient(ctx context.Context, req *pb.NotifyClientRequest) (*pb.Empty, error) {
	// Implement logic to handle notification to client
	// Step 1: Notify client about upload success or failure
	return &pb.Empty{}, nil
}

func (s *MasterTrackerServer) PingMasterTracker(ctx context.Context, req *pb.PingMasterTrackerRequest) (*pb.Empty, error) {
	// Implement logic to handle ping from data keeper node
	// Step 1: Update the lookup table loop through the lookup table and check if the data keeper node is alive
	for i := 0; i < len(s.lookupTable); i++ {
		if s.lookupTable[i].DataKeeperNode == req.NodeName {
			s.lookupTable[i].IsDataNodeAlive = true
			break
		}
	}
	return &pb.Empty{}, nil
}

func (s *MasterTrackerServer) UploadSuccess(ctx context.Context, req *pb.UploadSuccessRequest) (*pb.Empty, error) {
	// Implement logic to handle notification from data keeper node about successful upload
	// Step 1: Update the lookup table
	for i := 0; i < len(s.lookupTable); i++ {
		if s.lookupTable[i].DataKeeperNode == req.DataKeeperNodeName {
			s.lookupTable[i].FileName = req.FileName
			s.lookupTable[i].FilePath = req.FilePathOnNode
			break
		}
	}
	return &pb.Empty{}, nil
}

func main() {
	// Start gRPC server
	port := 8080
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port)) // Change port if needed
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	defer lis.Close()
	// Initialize MasterTrackerServer
	masterTracker := &MasterTrackerServer{
		lookupTable: []FileRecord{},
	}
	grpcServer := grpc.NewServer()
	pb.RegisterDFSServer(grpcServer, masterTracker)
	fmt.Println("Server started at port :", port) // Change port if needed

	if err := grpcServer.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}
