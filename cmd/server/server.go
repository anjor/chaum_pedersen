package main

import (
	"chaum_pedersen/zkp_auth"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	g int64 = 4
	h int64 = 9
	q int64 = 23
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	authServer := zkp_auth.NewServer("./server.db", g, h, q)
	zkp_auth.RegisterAuthServer(s, authServer)

	fmt.Println("gRPC server is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
