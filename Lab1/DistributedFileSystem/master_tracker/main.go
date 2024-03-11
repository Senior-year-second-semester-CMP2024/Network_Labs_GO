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
	// lookupTable map[string]FileRecord
}

func (s *MasterTrackerServer) RequestToUpload(ctx context.Context, req *pb.Empty) (*pb.RequestToUploadResponse, error) {
	// Implement logic to handle client request to upload a file
	// Step 1: Respond with a token for the client to use
	return &pb.RequestToUploadResponse{
		Token: "some_token",
	}, nil
}

func (s *MasterTrackerServer) UploadFile(ctx context.Context, req *pb.UploadFileRequest) (*pb.UploadFileResponse, error) {
	// Implement logic to handle file upload from client
	// Step 1: Receive the file data from the client
	// Step 2: Transfer the file to a data keeper node
	// Step 3: Update lookup table with file record
	// Step 4: Replicate file to 2 other nodes
	return &pb.UploadFileResponse{
		Message: "Success",
	}, nil
}

func (s *MasterTrackerServer) NotifyClient(ctx context.Context, req *pb.NotifyClientRequest) (*pb.Empty, error) {
	// Implement logic to handle notification to client
	// Step 1: Notify client about upload success or failure
	return &pb.Empty{}, nil
}

// Function to perform replication check
func replicationCheck(mt *MasterTrackerServer) {
	// for {
	// 	// Sleep for 10 seconds
	// 	time.Sleep(10 * time.Second)

	// 	// Replication algorithm
	// 	for _, record := range mt.lookupTable {
	// 		for getInstanceCount(record.FileName) < 3 {
	// 			sourceMachine := getSourceMachine(record)
	// 			destinationMachine := selectMachineToCopyTo()
	// 			notifyMachineDataTransfer(sourceMachine, destinationMachine, record)
	// 		}
	// 	}
	// }
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
	masterTracker := &MasterTrackerServer{}
	grpcServer := grpc.NewServer()
	pb.RegisterDFSServer(grpcServer, masterTracker)
	fmt.Println("Server started at port :", port) // Change port if needed

	if err := grpcServer.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v", err)
	}
	// Start replication check in a separate goroutine
	go replicationCheck(masterTracker)

	if err := grpcServer.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v", err)
	}
}
