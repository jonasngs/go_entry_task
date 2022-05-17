package services

import (
	"crypto/rand"
	"encoding/base64"
	"log"

	pb "github.com/jonasngs/go_entry_task/grpc"
)

type SessionInterface interface {
	VerifySession(sessionToken string) (*pb.User, error)
	GenerateSessionToken() string
	CreateSession(username string, sessionToken string) error
	// deleteSession()
}

type SessionService struct {
	cs  CacheInterface
	dao DAOInterface
}

func InitializeSession(cs CacheInterface, dao DAOInterface) SessionInterface {
	return SessionService{cs: cs, dao: dao}
}

// Create new user session and store in cache
func (sessionService SessionService) CreateSession(username string, sessionToken string) error {
	user, err := sessionService.dao.fetchUser(username)
	if err != nil {
		log.Printf("Unable to create session: %s\n", err)
		return err
	}
	sessionService.cs.updateCache(sessionToken, user)
	return nil
}

// Generate session token for session based authentication
func (sessionService SessionService) GenerateSessionToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("Unable to generate random byte sequence: %s\n", err)
	}
	return base64.URLEncoding.EncodeToString(b)
}

// Verifies if the user already has a session
func (sessionService SessionService) VerifySession(sessionToken string) (*pb.User, error) {
	user, err := sessionService.cs.checkCache(sessionToken)
	if err != nil {
		log.Printf("Could not retrieve user session: %s\n", err)
		return nil, err
	}
	return user, nil

}
