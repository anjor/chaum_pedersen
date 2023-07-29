package main

import (
	"chaum_pedersen/zkp_auth"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"math/rand"
)

const (
	g int64 = 4
	h int64 = 9
	q int64 = 23
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := zkp_auth.NewAuthClient(conn)

	// Register flow
	user := "test_user"
	var secret int64 = 1 // assumed to be mod q
	y1 := zkp_auth.Pow(g, secret)
	y2 := zkp_auth.Pow(h, secret)
	fmt.Printf("y1=%d, y2=%d\n", y1, y2)
	registerRequest := &zkp_auth.RegisterRequest{
		User: user,
		Y1:   y1,
		Y2:   y2,
	}
	registerResponse, err := client.Register(context.Background(), registerRequest)
	if err != nil {
		log.Fatalf("Register failed: %v", err)
	}
	fmt.Println("Register Response:", registerResponse)

	// Login flow

	// Commitment step
	k := zkp_auth.Mod(rand.Int63()+1, q)
	r1 := zkp_auth.Pow(g, k)
	r2 := zkp_auth.Pow(h, k)
	fmt.Printf("k=%d, r1=%d, r2=%d\n", k, r1, r2)
	challengeRequest := &zkp_auth.AuthenticationChallengeRequest{
		User: user,
		R1:   r1,
		R2:   r2,
	}
	challengeResponse, err := client.CreateAuthenticationChallenge(context.Background(), challengeRequest)
	if err != nil {
		log.Fatalf("CreateAuthenticationChallenge failed: %v", err)
	}
	fmt.Println("Authentication Challenge Response:", challengeResponse)

	// Response step
	s := zkp_auth.Mod(k-challengeResponse.C*secret, q)
	verifyRequest := &zkp_auth.AuthenticationAnswerRequest{
		AuthId: challengeResponse.GetAuthId(),
		S:      s,
	}
	fmt.Println("response: ", verifyRequest)
	verifyResponse, err := client.VerifyAuthentication(context.Background(), verifyRequest)
	if err != nil {
		log.Fatalf("VerifyAuthentication failed: %v", err)
	}
	fmt.Println("Authentication Answer Response:", verifyResponse)
}
