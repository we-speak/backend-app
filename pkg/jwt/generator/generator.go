package generator

import (
	"backend-app/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateTokenPair(userID uint, role string) (config.TokenPair, error) {
	accessClaims := &config.Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.AccessTokenExpiry)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(config.JWTSecret)
	if err != nil {
		return config.TokenPair{}, err
	}

	refreshClaims := &config.Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.RefreshTokenExpiry)),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(config.RefreshJWTSecret)
	if err != nil {
		return config.TokenPair{}, err
	}

	return config.TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}
