package main

import (
	"chaum_pedersen/zkp_auth"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := zkp_auth.NewAuthClient(conn)
	user := "test_user"
	secret := int64(10)

	// Register flow
	if err := zkp_auth.Register(client, user, secret); err != nil {
		log.Fatalf("Register failed: %v", err)
	}

	// Login flow
	sid, err := zkp_auth.Login(client, user, secret)
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	fmt.Printf("Session id = %s\n", sid)
}
