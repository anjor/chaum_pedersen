package main

import (
	"chaum_pedersen/zkp_auth"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	authServer := zkp_auth.NewServer("./server.db", 4, 9, 23)
	zkp_auth.RegisterAuthServer(s, authServer)

	fmt.Println("gRPC server is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
