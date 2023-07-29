package zkp_auth

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"math/rand"
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
	c := int64(rand.Intn(8))

	fmt.Printf("authId=%s, u=%s, r1=%d, r2=%d, c=%d\n\n", authId, req.GetUser(), req.GetR1(), req.GetR2(), c)
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
	expR1 := Pow(s.params.g, req.S) * Pow(y1, c)
	expR2 := Pow(s.params.h, req.S) * Pow(y2, c)

	fmt.Printf("r1 = %d, expR1=%d, r2=%d, expR2=%d\n", r1, expR1, r2, expR2)

	if expR1 == r1 && expR2 == r2 {
		sid := "session_" + uuid.New().String()
		return &AuthenticationAnswerResponse{SessionId: sid}, nil
	}
	return nil, fmt.Errorf("failed to verify authentication")
}

func NewServer(dbPath string, g, h, q int64) AuthServer {
	pp := protocolParams{
		g: g,
		h: h,
		q: q,
	}
	ls := make(map[uuid.UUID]*loginSession)
	return &authServer{dbPath: dbPath, loginSessions: ls, params: pp}
}
