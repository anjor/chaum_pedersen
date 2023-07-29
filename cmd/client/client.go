package main

import (
	"chaum_pedersen/zkp_auth"
	"context"
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

	registerRequest := &zkp_auth.RegisterRequest{
		User: "test_user",
		Y1:   1,
		Y2:   2,
	}

	registerResponse, err := client.Register(context.Background(), registerRequest)
	if err != nil {
		log.Fatalf("Register failed: %v", err)
	}
	fmt.Println("Register Response:", registerResponse)

	// Call the CreateAuthenticationChallenge RPC
	challengeRequest := &zkp_auth.AuthenticationChallengeRequest{
		User: "test_user",
		R1:   789,
		R2:   1011,
	}
	challengeResponse, err := client.CreateAuthenticationChallenge(context.Background(), challengeRequest)
	if err != nil {
		log.Fatalf("CreateAuthenticationChallenge failed: %v", err)
	}
	fmt.Println("Authentication Challenge Response:", challengeResponse)

	// Call the VerifyAuthentication RPC
	verifyRequest := &zkp_auth.AuthenticationAnswerRequest{
		AuthId: challengeResponse.GetAuthId(),
		S:      1213,
	}
	verifyResponse, err := client.VerifyAuthentication(context.Background(), verifyRequest)
	if err != nil {
		log.Fatalf("VerifyAuthentication failed: %v", err)
	}
	fmt.Println("Authentication Answer Response:", verifyResponse)
}
