package auth

import (
	"context"
	"errors"
	"time"

	"github.com/mtsfy/fotosouk/internal/models"
	"github.com/mtsfy/fotosouk/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo UserRepository
}

func NewAuthService(repo UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(ctx context.Context, firstName, lastName, email, username, password string) (*models.User, error) {
	if len(username) == 0 {
		return nil, errors.New("username is required")
	}

	if len(email) == 0 {
		return nil, errors.New("email is required")
	}

	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	exists, err := s.repo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	exists, err = s.repo.ExistsByUserName(ctx, username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already taken")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &models.User{
		FirstName:      firstName,
		LastName:       lastName,
		Email:          email,
		Username:       username,
		HashedPassword: string(hash),
	}

	return s.repo.Create(ctx, u)
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*models.User, *utils.Token, error) {
	if len(username) == 0 {
		return nil, nil, errors.New("username is required")
	}

	if len(password) == 0 {
		return nil, nil, errors.New("password is required")
	}

	u, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, nil, errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	if err != nil {
		return nil, nil, errors.New("invalid username or password")
	}

	token, err := utils.GenerateToken(u.ID, u.Username)
	if err != nil {
		return nil, nil, err
	}

	rt := &models.RefreshToken{
		UserID:    u.ID,
		Token:     token.RefreshToken,
		ExpiresAt: time.Unix(token.RefreshExpiresAt, 0),
	}
	if err := s.repo.SaveRefreshToken(ctx, rt); err != nil {
		return nil, nil, err
	}

	return u, token, nil
}

func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (*utils.Token, error) {
	// check refresh token exists and is not revoked
	rt, err := s.repo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// check expiration
	if time.Now().After(rt.ExpiresAt) {
		return nil, errors.New("refresh token expired")
	}

	// check jwt signature
	claims, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	userID := int(claims["user_id"].(float64))
	username := claims["username"].(string)

	// create new access token
	token, err := utils.GenerateToken(userID, username)
	if err != nil {
		return nil, err
	}

	return token, nil
}
