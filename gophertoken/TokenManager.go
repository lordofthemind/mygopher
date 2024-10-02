package gophertoken

import (
	"time"

	"github.com/google/uuid"
)

const (
	TokenTypeJWT    = "jwt"
	TokenTypePaseto = "paseto"
)

// TokenManager is the interface for creating and verifying tokens.
//
// Example usage:
//
//	var manager TokenManager
//	manager, err = NewTokenManager("jwt", "your-secret-key")
//	if err != nil {
//	  log.Fatal(err)
//	}
type TokenManager interface {
	GenerateToken(userID uuid.UUID, username string, duration time.Duration) (string, error)
	ValidateToken(token string) (*Payload, error)
}

// NewTokenManager creates a new token manager (JWT or Paseto) depending on the provided type.
//
// Example usage:
//
//	manager, err := NewTokenManager("jwt", "your-secret-key")
//	if err != nil {
//	  log.Fatal(err)
//	}
func NewTokenManager(tokenType, secretKey string) (TokenManager, error) {
	switch tokenType {
	case TokenTypeJWT:
		return NewJWTMaker(secretKey)
	case TokenTypePaseto:
		return NewPasetoMaker(secretKey)
	default:
		return nil, ErrInvalidToken
	}
}
