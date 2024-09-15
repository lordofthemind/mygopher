package gophertoken

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTMaker is a struct for handling JWT token creation and validation.
type JWTMaker struct {
	symmetricKey string
}

// NewJWTMaker creates a new JWTMaker with the given symmetric key.
//
// Example usage:
//
//	maker, err := NewJWTMaker("your-secret-key")
//	if err != nil {
//	  log.Fatal(err)
//	}
func NewJWTMaker(secretKey string) (TokenManager, error) {
	if len(secretKey) == 0 {
		return nil, errors.New("symmetric key must be set")
	}
	return &JWTMaker{symmetricKey: secretKey}, nil
}

// GenerateToken creates a new JWT token for a specific user with a given duration.
//
// Example usage:
//
//	token, err := maker.GenerateToken("user123", time.Hour)
//	if err != nil {
//	  log.Fatal(err)
//	}
func (j *JWTMaker) GenerateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"id":         payload.ID.String(),
		"username":   payload.Username,
		"issued_at":  payload.IssuedAt.Unix(),
		"expired_at": payload.ExpiredAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.symmetricKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken checks if the given JWT token is valid.
//
// Example usage:
//
//	payload, err := maker.ValidateToken(tokenString)
//	if err != nil {
//	  log.Fatal("Invalid token")
//	}
func (j *JWTMaker) ValidateToken(tokenString string) (*Payload, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.symmetricKey), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	payload := &Payload{
		ID:        uuid.MustParse(claims["id"].(string)),
		Username:  claims["username"].(string),
		IssuedAt:  time.Unix(int64(claims["issued_at"].(float64)), 0),
		ExpiredAt: time.Unix(int64(claims["expired_at"].(float64)), 0),
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
