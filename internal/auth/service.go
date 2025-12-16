package auth

import (
	"context"
	"errors"

	"github.com/mtsfy/fotosouk/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo UserRepository
}

func NewAuthService(repo UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(ctx context.Context, firstName, lastName, email, username, password string) (*user.User, error) {
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

	u := &user.User{
		FirstName:      firstName,
		LastName:       lastName,
		Email:          email,
		Username:       username,
		HashedPassword: string(hash),
	}

	return s.repo.Create(ctx, u)
}
