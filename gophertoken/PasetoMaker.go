package gophertoken

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

// PasetoMaker is a struct for handling Paseto token creation and validation.
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// NewPasetoMaker creates a new PasetoMaker with the given symmetric key.
//
// Example usage:
//
//	maker, err := NewPasetoMaker("your-secret-key")
//	if err != nil {
//	  log.Fatal(err)
//	}
func NewPasetoMaker(secretKey string) (TokenManager, error) {
	if len(secretKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d bytes", chacha20poly1305.KeySize)
	}
	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(secretKey),
	}
	return maker, nil
}

// GenerateToken creates a new Paseto token for a specific user with a given duration.
//
// Example usage:
//
//	token, err := maker.GenerateToken("user123", time.Hour)
//	if err != nil {
//	  log.Fatal(err)
//	}
func (maker *PasetoMaker) GenerateToken(userID uuid.UUID, username string, duration time.Duration) (string, error) {
	// Create the payload with userID and username
	payload, err := NewPayload(userID, username, duration)
	if err != nil {
		return "", err
	}

	// Encrypt the payload and return the token string
	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
}

// ValidateToken checks if the given Paseto token is valid.
//
// Example usage:
//
//	payload, err := maker.ValidateToken(tokenString)
//	if err != nil {
//	  log.Fatal("Invalid token")
//	}
func (maker *PasetoMaker) ValidateToken(token string) (*Payload, error) {
	// Decrypt the token to extract the payload
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Validate the payload (check expiration)
	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
