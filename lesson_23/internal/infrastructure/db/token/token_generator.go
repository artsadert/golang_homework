package token

import (
	"os"

	"github.com/go-chi/jwtauth/v5"
)

func GetAccessToken() *jwtauth.JWTAuth {
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET_KEY")), nil) // replace with secret key

	return tokenAuth
}

func GetOldAccessToken() *jwtauth.JWTAuth {
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_OLD_SECRET_KEY")), nil) // replace with secret key

	return tokenAuth
}

func GetRefreshToken() *jwtauth.JWTAuth {
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_REFRESH_SECRET_KEY")), nil) // replace with secret key

	return tokenAuth
}
