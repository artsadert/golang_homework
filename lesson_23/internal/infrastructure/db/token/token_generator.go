package token

import (
	"os"

	"github.com/go-chi/jwtauth/v5"
)

func GetToken() *jwtauth.JWTAuth {
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET_KEY")), nil) // replace with secret key

	return tokenAuth
}
