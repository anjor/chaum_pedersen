package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
)

type loginSession struct {
	u  string
	r1 int64
	r2 int64
	c  int64
}

type protocolParams struct {
	g int64
	h int64
	q int64
}

type authServer struct {
	dbPath        string
	loginSessions map[uuid.UUID]*loginSession // map from auth id to login session
	params        protocolParams
}

func (s *authServer) getProtocolparams(ctx context.Context) *protocolParams {
	return &s.params
}

func (s *authServer) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	err := insertRowIntoDB(s.dbPath, req.GetUser(), req.GetY1(), req.GetY2())
	if err != nil {
		return nil, err
	}
	return &RegisterResponse{}, nil
}

func (s *authServer) CreateAuthenticationChallenge(ctx context.Context, req *AuthenticationChallengeRequest) (*AuthenticationChallengeResponse, error) {
	// generate auth id
	authId := uuid.New()
	// generate challenge c
	c := rand.Int63()

	// persist user, r1, r2 and c in loginSession
	s.loginSessions[authId] = &loginSession{
		u:  req.GetUser(),
		r1: req.GetR1(),
		r2: req.GetR2(),
		c:  c,
	}

	// construct response and send it back
	response := &AuthenticationChallengeResponse{
		AuthId: authId.String(),
		C:      c,
	}
	return response, nil
}

func (s *authServer) VerifyAuthentication(ctx context.Context, req *AuthenticationAnswerRequest) (*AuthenticationAnswerResponse, error) {
	// look up login session from auth id
	authId, err := uuid.Parse(req.AuthId)
	if err != nil {
		return nil, fmt.Errorf("invalid auth id. Expected UUID, instead got %s. err=%s\n", req.AuthId, err)
	}

	ls := s.loginSessions[authId]
	u := ls.u
	r1 := ls.r1
	r2 := ls.r2
	c := ls.c

	// fetch y1 and y2 from database based on user
	_, y1, y2, err := getRowByUser(s.dbPath, u)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve registration y1 and y2. err=%s\n", err)
	}

	// calculate g^s*y1^c and h^s*y2^c
	expR1 := pow(s.params.g, req.S) * pow(y1, c)
	expR2 := pow(s.params.h, req.S) * pow(y2, c)

	if expR1 == r1 && expR2 == r2 {
		sid := "session_" + uuid.New().String()
		return &AuthenticationAnswerResponse{SessionId: sid}, nil
	}
	return nil, fmt.Errorf("failed to verify authentication")
}

func insertRowIntoDB(dbPath string, user string, y1 int64, y2 int64) error {
	// Open the database file
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Create the table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS user_table (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			u TEXT,
			y1 INTEGER,
			y2 INTEGER
		);`)
	if err != nil {
		return err
	}

	// Insert the data into the table
	_, err = db.Exec("INSERT INTO my_table (u, y1, y2) VALUES (?, ?, ?)",
		user, y1, y2)
	if err != nil {
		return err
	}

	fmt.Printf("Row inserted successfully for user: %s\n", user)
	return nil
}

func getRowByUser(dbPath, user string) (string, int64, int64, error) {
	// Open the database file
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return "", 0, 0, err
	}
	defer db.Close()

	// Query the row with the specified user
	row := db.QueryRow("SELECT u, y1, y2 FROM user_table WHERE u = ?", user)

	// Extract the data from the row
	var retrievedUser string
	var y1, y2 int64
	err = row.Scan(&retrievedUser, &y1, &y2)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", 0, 0, fmt.Errorf("user not found: %s", user)
		}
		return "", 0, 0, err
	}

	return retrievedUser, y1, y2, nil
}

// See Donald Knuth, The Art of Computer Programming, Volume 2, Section 4.6.3
func pow(a, b int64) int64 {
	var p int64 = 1
	for b > 0 {
		if b&1 != 0 {
			p *= a
		}
		b >>= 1
		a *= a
	}
	return p
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterAuthServer(s, &authServer{})
	fmt.Println("gRPC server is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
