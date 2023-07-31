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

	if err := success(client, user, secret); err != nil {
		fmt.Printf("register and login should have succeded. Instead failed with error: %v\n", err)
	}

	if err := fail(client, user, secret); err == nil {
		fmt.Printf("register and login should have failed. Instead succeded")
	}

}

func success(client zkp_auth.AuthClient, user string, secret int64) error {
	// Register flow
	if err := zkp_auth.Register(client, user, secret); err != nil {
		return err
	}
	// Login flow
	sid, err := zkp_auth.Login(client, user, secret)
	if err != nil {
		return err
	}

	fmt.Printf("session id = %s\n", sid)
	return nil
}

func fail(client zkp_auth.AuthClient, user string, secret int64) error {
	// Register flow
	if err := zkp_auth.Register(client, user, secret); err != nil {
		return err
	}
	// Login flow
	sid, err := zkp_auth.Login(client, user, secret+1)
	if err != nil {
		return err
	}

	fmt.Printf("Session id = %s\n", sid)
	return nil

}
