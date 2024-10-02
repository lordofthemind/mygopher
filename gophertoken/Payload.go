package gophertoken

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Errors related to token validation.
var (
	ErrInvalidToken = errors.New("token validation failed: signature invalid or claims malformed")
	ErrExpiredToken = errors.New("token validation failed: token has expired")
)

// Payload contains the data embedded within a token.
type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates a new token payload with a specific username and token duration.
//
// Example usage:
//
//	payload, err := NewPayload(userID, "username123", time.Hour)
//	if err != nil {
//	  log.Fatal(err)
//	}
func NewPayload(userID uuid.UUID, username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		UserID:    userID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

// Valid checks if the payload's expiration date has passed and returns an error if it has.
//
// Example usage:
//
//	err := payload.Valid()
//	if err != nil {
//	  log.Fatal("Token expired")
//	}
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
