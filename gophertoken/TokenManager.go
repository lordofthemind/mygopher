package tokens

import (
	"time"

	"github.com/lordofthemind/htmx_GO/internals/configs"
)

type TokenManager interface {
	// CreateToken creates a new token for a specific user
	GenerateToken(username string, duration time.Duration) (string, error)
	// VerifyToken checks if the token is valid or not
	ValidateToken(tokenString string) (*Payload, error)
}

func NewTokenManager() (TokenManager, error) {
	if configs.UseJWT {
		return NewJWTMaker()
	}
	return NewPasetoMaker()
}
