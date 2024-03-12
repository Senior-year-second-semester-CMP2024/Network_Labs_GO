package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "DistributedFileSystem/dfs" // Import generated proto file

	"google.golang.org/grpc"
)

type FileRecord struct {
	FileName        []string
	Ports           []string
	FilePath        []string
	IsDataNodeAlive bool
}

// MasterTrackerServer implements the DFS service
type MasterTrackerServer struct {
	pb.UnimplementedDFSServer
	client      pb.DFSClient
	lookupTable map[string]FileRecord // Lookup table to store the file records key = DataKeeperNode
}

func (s *MasterTrackerServer) RequestToUpload(ctx context.Context, req *pb.Empty) (*pb.RequestToUploadResponse, error) {
	// Implement logic to handle client request to upload a file
	// token = port number of the data keeper node that exist in the lookup table
	token := "50051"
	log.Println("Request to Upload")
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
	if _, ok := s.lookupTable[req.NodeName]; !ok {
		// If not, create a new FileRecord for this node
		s.lookupTable[req.NodeName] = FileRecord{
			FileName:        []string{},
			Ports:           []string{},
			FilePath:        []string{},
			IsDataNodeAlive: true,
		}
	}
	record := s.lookupTable[req.NodeName]
	record.IsDataNodeAlive = true
	s.lookupTable[req.NodeName] = record

	log.Println("Data node is alive:", req.NodeName, " port : ")
	return &pb.Empty{}, nil
}

func (s *MasterTrackerServer) UploadSuccess(ctx context.Context, req *pb.UploadSuccessRequest) (*pb.Empty, error) {
	// Implement logic to handle notification from data keeper node about successful upload
	// Step 1: Update the lookup table
	// Check if the lookupTable already contains the DataKeeperNodeName
	if _, ok := s.lookupTable[req.DataKeeperNodeName]; !ok {
		// If not, create a new FileRecord for this node
		s.lookupTable[req.DataKeeperNodeName] = FileRecord{
			FileName:        []string{},
			Ports:           []string{},
			FilePath:        []string{},
			IsDataNodeAlive: true,
		}
	}

	nodeRecord := s.lookupTable[req.DataKeeperNodeName]
	nodeRecord.FileName = append(nodeRecord.FileName, req.FileName)
	nodeRecord.Ports = append(nodeRecord.Ports, req.Ports)
	nodeRecord.FilePath = append(nodeRecord.FilePath, req.FilePathOnNode)

	s.lookupTable[req.DataKeeperNodeName] = nodeRecord

	// Call the UploadSuccess RPC using the client
	// Prepare the request
	request := &pb.NotifyClientRequest{
		Message: "success",
	}
	_, err := s.client.NotifyClient(context.Background(), request)
	log.Println("Upload success:", req.FileName)
	if err != nil {
		return &pb.Empty{}, err
	}
	return &pb.Empty{}, nil
}

func main() {
	// Client setup
	// Set up a gRPC connection to the server implementing UploadSuccess
	ClientConn, err := grpc.Dial("localhost:8081", grpc.WithInsecure()) // Update with actual server address
	if err != nil {
		log.Fatalf("failed to connect to client: %v", err)
	}
	defer ClientConn.Close()
	client := pb.NewDFSClient(ClientConn)

	// Start Master gRPC server
	port := ":8080"
	lis, err := net.Listen("tcp", port) // Change port if needed
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	defer lis.Close()
	// Initialize MasterTrackerServer
	masterTracker := &MasterTrackerServer{
		lookupTable: make(map[string]FileRecord),
		client:      client,
	}
	grpcServer := grpc.NewServer()
	pb.RegisterDFSServer(grpcServer, masterTracker)
	fmt.Println("Master server started at ", port) // Change port if needed

	if err := grpcServer.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}
