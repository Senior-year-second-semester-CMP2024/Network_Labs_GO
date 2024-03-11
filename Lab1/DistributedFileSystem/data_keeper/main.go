package main

import (
	"fmt"
	// "context"
	// "net"
	// "strings"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		fmt.Println("did not connect:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Hello, world!")
}
