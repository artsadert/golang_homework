package user

import (
	"time"

	"github.com/artsadert/lesson_23/internal/domain/entities"
	"github.com/artsadert/lesson_23/internal/interface/api/rest/v1/user/dto/response"
	"github.com/google/uuid"
)

func createTokens(uuid uuid.UUID, config *entities.Config) (*response.LoginResponse, error) {
	claims_access := map[string]interface{}{
		"type":      "access",
		"exp":       time.Now().Add(15 * time.Minute),
		"iat":       time.Now(),
		"user_uuid": uuid.String(),
	}

	claims_refresh := map[string]interface{}{
		"type": "refresh",
		"iat":  time.Now(),
		"exp":  time.Now().Add(24 * time.Hour),
	}

	_, accessTokenString, err := config.JWTAccessSecret.Encode(claims_access)
	if err != nil {
		return nil, err
	}

	_, refreshTokenString, err := config.JWTRefreshSecret.Encode(claims_refresh)
	if err != nil {
		return nil, err
	}

	return &response.LoginResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}
