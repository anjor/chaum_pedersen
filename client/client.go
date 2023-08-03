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
	user := "test_user_1"
	secret := int64(10)

	fmt.Println("Successful attempt started")
	if err := success(client, user, secret); err != nil {
		fmt.Printf("register and login should have succeded. Instead failed with error: %v\n", err)
	} else {
		fmt.Println("Test for successful login attempt started succeeded")
	}
	fmt.Println()

	user = "test_user_2"
	secret = int64(100)

	fmt.Println("Failed attempt started")
	if err := fail(client, user, secret); err == nil {
		fmt.Printf("register and login should have failed. Instead succeded")
	} else {
		fmt.Println("Test for failed login attempt succeeded")
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
