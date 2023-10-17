package http

import (
	"os"

	"github.com/go-chi/jwtauth/v5"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	privateKey := os.Getenv("JWT_PRIVATE_KEY")
	if privateKey == "" {
		privateKey = "secret"
	}

	tokenAuth = jwtauth.New("HS256", []byte(privateKey), nil)
}
