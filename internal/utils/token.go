package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mtsfy/fotosouk/internal/config"
)

type Token struct {
	AccessToken     string `json:"accessToken"`
	AccessExpiresAt int64  `json:"accessExpiresAt"`

	RefreshToken     string `json:"refreshToken"`
	RefreshExpiresAt int64  `json:"refreshExpiresAt"`
}

func GenerateToken(id int, userName string) (*Token, error) {
	t := &Token{}
	var err error

	accessSecret := config.Config("JWT_ACCESS_SECRET")
	if accessSecret == "" {
		return nil, errors.New("JWT_ACCESS_SECRET is not set")
	}

	refreshSecret := config.Config("JWT_REFRESH_SECRET")
	if refreshSecret == "" {
		return nil, errors.New("JWT_REFRESH_SECRET is not set")
	}

	t.AccessExpiresAt = time.Now().Add(24 * time.Hour).Unix() // 1 day
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  id,
		"username": userName,
		"exp":      t.AccessExpiresAt,
	})

	t.AccessToken, err = accessToken.SignedString([]byte(accessSecret))
	if err != nil {
		return nil, err
	}

	t.RefreshExpiresAt = time.Now().Add(24 * time.Hour * 7).Unix() // 7 days
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  id,
		"username": userName,
		"exp":      t.RefreshExpiresAt,
	})

	t.RefreshToken, err = refreshToken.SignedString([]byte(refreshSecret))
	if err != nil {
		return nil, err
	}

	return t, nil
}

func ValidateRefreshToken(tokenString string) (jwt.MapClaims, error) {
	refreshSecret := config.Config("JWT_REFRESH_SECRET")
	if refreshSecret == "" {
		return nil, errors.New("JWT_REFRESH_SECRET is not set")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(refreshSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
