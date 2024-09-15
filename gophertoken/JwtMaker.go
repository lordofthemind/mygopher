package tokens

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lordofthemind/htmx_GO/internals/configs"
)

type JWTMaker struct{}

// NewJWTMaker creates a new JWTMaker
func NewJWTMaker() (TokenManager, error) {
	if len(configs.TokenSymmetricKey) == 0 {
		return nil, errors.New("symmetric key must be set in the configuration")
	}
	return &JWTMaker{}, nil
}

// GenerateToken creates a new token for a specific user
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

	tokenString, err := token.SignedString([]byte(configs.TokenSymmetricKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken checks if the token is valid or not
func (j *JWTMaker) ValidateToken(tokenString string) (*Payload, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(configs.TokenSymmetricKey), nil
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
