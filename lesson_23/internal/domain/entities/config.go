package entities

import "github.com/go-chi/jwtauth/v5"

type Config struct {
	JWTAccessSecret    *jwtauth.JWTAuth
	JWTRefreshSecret   *jwtauth.JWTAuth
	JWTOLDAccessSecret *jwtauth.JWTAuth
}
