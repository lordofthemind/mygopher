package tokens

import (
	"fmt"
	"time"

	"github.com/lordofthemind/htmx_GO/internals/configs"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

// PasetoMaker is a struct that holds the paseto instance and the symmetric key
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// NewPasetoMaker creates a new PasetoMaker
func NewPasetoMaker() (TokenManager, error) {
	if len(configs.TokenSymmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid secret key size: must be exactly %d bytes", chacha20poly1305.KeySize)
	}
	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(configs.TokenSymmetricKey),
	}
	return maker, nil
}

// GenerateToken creates a new token for a specific user
func (maker *PasetoMaker) GenerateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
}

// ValidateToken checks if the token is valid or not
func (maker *PasetoMaker) ValidateToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
