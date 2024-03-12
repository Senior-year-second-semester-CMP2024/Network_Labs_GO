package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	pb "DistributedFileSystem/dfs" // Import generated proto file

	"google.golang.org/grpc"
)

type FileRecord struct {
	NodeName        string
	FileName        []string
	Ports           []string
	FilePath        []string
	IsDataNodeAlive bool
}

// MasterTrackerServer implements the DFS service
type MasterTrackerServer struct {
	pb.UnimplementedDFSServer
	client           pb.DFSClient
	lookupTable      map[string]FileRecord // Lookup table to store the file records key = DataKeeperNode
	distinctFilesSet Set
	mu               sync.Mutex
}

func (s *MasterTrackerServer) RequestToUpload(ctx context.Context, req *pb.Empty) (*pb.RequestToUploadResponse, error) {
	// Implement logic to handle client request to upload a file
	// token = port number of the data keeper node that exist in the lookup table

	// TODO : choose and alive node and return the port
	token := "6200" // initially
	keys := make([]string, 0, len(s.lookupTable))
	for k := range s.lookupTable {
		keys = append(keys, k)
	}
	for { // loop until an alive node is found
		randNode := keys[rand.Intn(len(keys))]
		dataNode := s.lookupTable[randNode] // select a random alive node
		if dataNode.IsDataNodeAlive {
			randPort := rand.Intn(len(s.lookupTable[randNode].Ports))
			token = dataNode.Ports[randPort] // select a random port
			log.Println("Request to Upload Node: '", randNode, "' on Port:'", token, "'")
			return &pb.RequestToUploadResponse{
				Token: token,
			}, nil
		}
	}
}

func (s *MasterTrackerServer) PingMasterTracker(ctx context.Context, req *pb.PingMasterTrackerRequest) (*pb.Empty, error) {
	// Lock the mutex before accessing the shared data structures
	s.mu.Lock()
	defer s.mu.Unlock()

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
	record.Ports = req.AvailablePorts
	record.NodeName = req.NodeName
	s.lookupTable[req.NodeName] = record

	// log.Println("Data node is alive:", req.NodeName, " ports : ", s.lookupTable[req.NodeName].Ports, " lookuptable size : ", len(s.lookupTable))
	return &pb.Empty{}, nil
}

func (s *MasterTrackerServer) UploadSuccess(ctx context.Context, req *pb.UploadSuccessRequest) (*pb.Empty, error) {
	// Implement logic to handle notification from data keeper node about successful upload
	// Step 1: Update the lookup table
	// Check if the lookupTable already contains the DataKeeperNodeName
	nodeRecord := s.lookupTable[req.DataKeeperNodeName]
	nodeRecord.FileName = append(nodeRecord.FileName, req.FileName)
	nodeRecord.FilePath = append(nodeRecord.FilePath, req.FilePathOnNode)

	s.lookupTable[req.DataKeeperNodeName] = nodeRecord
	s.distinctFilesSet.Add(req.FileName)

	// Call the UploadSuccess RPC using the client
	// Prepare the request
	request := &pb.NotifyClientRequest{
		Message: "success",
	}
	_, err := s.client.NotifyClient(context.Background(), request)
	log.Println("Upload success:", req.FileName, " on node:", req.DataKeeperNodeName)
	if err != nil {
		return &pb.Empty{}, err
	}
	return &pb.Empty{}, nil
}

func (s *MasterTrackerServer) RequestToDownload(ctx context.Context, req *pb.RequestToDownloadRequest) (*pb.RequestToDownloadResponse, error) {
	fileName := req.FileName
	var machineInfos []*pb.MachineInfo
	// seach on the lookup table for the file
	// for each data_node in the lookup table
	log.Print(s.lookupTable)
	for _, data_node := range s.lookupTable {
		// for each file in the data_node
		for _, name := range data_node.FileName {
			if name == fileName {
				// Assuming each data_node has corresponding port and filepath at index 'i'
				for _, port := range data_node.Ports {
					machineInfo := &pb.MachineInfo{}
					machineInfo.Port = port
					machineInfos = append(machineInfos, machineInfo)
				}
				log.Print("File found at node with port : ", data_node.Ports)
				return &pb.RequestToDownloadResponse{
					MachineInfos: machineInfos,
				}, nil
			}
		}
	}
	log.Print("File not found")
	return &pb.RequestToDownloadResponse{}, nil
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
	port := "8080"
	masterTracker := &MasterTrackerServer{
		lookupTable:      make(map[string]FileRecord),
		client:           client,
		distinctFilesSet: make(Set),
	}
	// Initialize MasterTrackerServer
	var wg sync.WaitGroup
	wg.Add(1)
	go startServer(port, &wg, masterTracker)

	// Start the repication routine
	go ReplicateRoutine(masterTracker)
	wg.Wait()
}

func startServer(port string, wg *sync.WaitGroup, masterTracker *MasterTrackerServer) {
	defer wg.Done()

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterDFSServer(grpcServer, masterTracker)
	fmt.Println("Master server started at ", port) // Change port if needed
	if err := grpcServer.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}

func ReplicateRoutine(s *MasterTrackerServer) {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for range ticker.C {
		// 1 - get distinct file instances from the lookup table
		// 2 - get source machine that has the file
		// 3 - while there are less than 3 instances of the file
		//     3.1 - get a random machine from the lookup table to replicate the file to
		//     3.2 - notify the source machine to replicate the file to the random machine
		distinctFiles := s.distinctFilesSet.ToList()
		randomMachinePort := ""
		randomMachine := ""
		for _, file := range distinctFiles {
			// 1
			sourceMachines, sourceMachinePort := GetSourceMachines(s.lookupTable, file)
			// 2
			// While loop instead of if
			for len(sourceMachines) < 3 {
				// 3
				// 3.1
				randomMachine, randomMachinePort = s.selectMachineToCopyTo(s.lookupTable, sourceMachines)
				NotifyMachine(file, sourceMachinePort, randomMachinePort)
				sourceMachines.Add(randomMachine)
				// 3.2
				log.Println("Lookup table : ", s.lookupTable)
			}
		}
	}
}

func GetSourceMachines(lookupTable map[string]FileRecord, fileName string) (Set, string) {
	sourceMachines := make(Set)
	sourceMachinePort := ""
	for _, data_node := range lookupTable { // search in each data node
		for _, name := range data_node.FileName { // search in each file
			if name == fileName { // if the file is found
				sourceMachines.Add(data_node.NodeName)
				// choose random number from 0 to len(ports - 1)
				rand := rand.Intn(len(lookupTable[data_node.NodeName].Ports))
				sourceMachinePort = lookupTable[data_node.NodeName].Ports[rand]
			}
		}
	}
	return sourceMachines, sourceMachinePort
}

func (s *MasterTrackerServer) selectMachineToCopyTo(lookupTable map[string]FileRecord, sourceMachines Set) (string, string) {
	// select a machine from the lookupTable that is not in the sourceMachines
	randomMachine := ""
	randomMachinePort := ""
	for data_node := range lookupTable {
		if !sourceMachines.Contains(data_node) {
			randomMachine = data_node
			// choose random number from 0 to len(ports - 1)
			rand := rand.Intn(len(lookupTable[data_node].Ports))
			randomMachinePort = lookupTable[data_node].Ports[rand]
		}
	}
	return randomMachine, randomMachinePort
}

func NotifyMachine(fileName string, sourcePort string, randomMachinePort string) {
	// TODO : Port of a node that has the file
	request := &pb.NotifyMachineDataTransferRequest{
		Filename: fileName,
		SrcPort:  sourcePort,
		DstPort:  randomMachinePort,
	}
	// connect to the destination port
	dataConn, err := grpc.Dial("localhost:"+sourcePort, grpc.WithInsecure()) // Update with actual server address
	if err != nil {
		log.Fatalf("failed to connect to data keeper: %v", err)
	}
	defer dataConn.Close()
	cData := pb.NewDFSClient(dataConn)
	// Call the UploadFile RPC
	_, err = cData.NotifyMachineDataTransfer(context.Background(), request)
	if err != nil {
		log.Println("Failed to call UploadFile:", err)
	}
}

// ---------------------------------------------------------------------//
// ------------------------ Set implementation--------------------------//
// ---------------------------------------------------------------------//
type Set map[string]bool

func (set Set) Add(element string) {
	set[element] = true
}

func (set Set) Remove(element string) {
	delete(set, element)
}

func (set Set) Contains(element string) bool {
	return set[element]
}

func (set Set) ToList() []string {
	list := make([]string, 0, len(set))
	for k := range set {
		list = append(list, k)
	}
	return list
}
