package entities

import (
	"time"

	"github.com/go-chi/jwtauth/v5"
)

type Config struct {
	JWTAccessSecret    *jwtauth.JWTAuth
	JWTRefreshSecret   *jwtauth.JWTAuth
	JWTOLDAccessSecret *jwtauth.JWTAuth

	JWT_TTL         time.Duration
	JWT_REFRESH_TTL time.Duration
}
