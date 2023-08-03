package zkp_auth

import (
	"context"
	"fmt"
	"log"
)

func Register(client AuthClient, user string, secret int64) error {
	y1, y2 := calculateCommitment(secret)
	fmt.Printf("registration: y1=%d, y2=%d\n", y1, y2)
	registerRequest := &RegisterRequest{
		User: user,
		Y1:   y1,
		Y2:   y2,
	}
	_, err := client.Register(context.Background(), registerRequest)
	return err
}

func Login(client AuthClient, user string, secret int64) (string, error) {
	// Commitment step
	k := generateRandom()

	cr, err := commit(client, user, k)
	if err != nil {
		log.Fatalf("commitment step failed: %v", err)
	}
	fmt.Println("authentication challenge response:", cr)

	// Response step
	return respond(client, cr.GetAuthId(), k, cr.GetC(), secret)
}

func commit(client AuthClient, user string, k int64) (*AuthenticationChallengeResponse, error) {
	r1, r2 := calculateCommitment(k)
	fmt.Printf("commitment: k=%d, r1=%d, r2=%d\n", k, r1, r2)
	challengeRequest := &AuthenticationChallengeRequest{
		User: user,
		R1:   r1,
		R2:   r2,
	}
	return client.CreateAuthenticationChallenge(context.Background(), challengeRequest)
}

func respond(client AuthClient, authId string, k, c, secret int64) (string, error) {
	s := mod(k-c*secret, q)
	verifyRequest := &AuthenticationAnswerRequest{
		AuthId: authId,
		S:      s,
	}
	fmt.Println("verify request: ", verifyRequest)
	verifyResponse, err := client.VerifyAuthentication(context.Background(), verifyRequest)
	if err != nil {
		return "", fmt.Errorf("respond step failed: %v", err)
	}
	return verifyResponse.GetSessionId(), nil
}
