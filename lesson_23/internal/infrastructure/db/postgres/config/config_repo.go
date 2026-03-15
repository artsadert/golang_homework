package config

import (
	"os"
	"time"

	"github.com/artsadert/lesson_23/internal/domain/entities"
	"github.com/artsadert/lesson_23/internal/domain/repository"
	"github.com/go-chi/jwtauth/v5"
)

type ConfigRepo struct{}

func NewConfigRepo() repository.ConfigRepo {
	return &ConfigRepo{}
}

func (c *ConfigRepo) GetConfig() (*entities.Config, error) {
	return &entities.Config{
		JWTAccessSecret:    c.getAccessToken(),
		JWTRefreshSecret:   c.getRefreshToken(),
		JWTOLDAccessSecret: c.getOldAccessToken(),
		JWT_TTL:            60 * time.Minute,
		JWT_REFRESH_TTL:    24 * time.Hour,
	}, nil
}

func (c *ConfigRepo) getAccessToken() *jwtauth.JWTAuth {
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET_KEY")), nil) // replace with secret key

	return tokenAuth
}

func (c *ConfigRepo) getOldAccessToken() *jwtauth.JWTAuth {
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_OLD_SECRET_KEY")), nil) // replace with secret key

	return tokenAuth
}

func (c *ConfigRepo) getRefreshToken() *jwtauth.JWTAuth {
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_REFRESH_SECRET_KEY")), nil) // replace with secret key

	return tokenAuth
}
