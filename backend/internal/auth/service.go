package auth

import (
	"context"
	"log/slog"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *UserRepo
	logger   *slog.Logger
}

func NewUserService(ur *UserRepo, logger *slog.Logger) *UserService {
	ts := UserService{
		userRepo: ur,
		logger:   logger,
	}
	return &ts
}

func (s *UserService) GetUserById(ctx context.Context, userId int) (*User, error) {
	return s.userRepo.getUserById(ctx, userId)
}

func (s *UserService) UsernameInUse(ctx context.Context, username string) (bool, error) {
	return s.userRepo.usernameInUse(ctx, strings.ToLower(username))
}

func (s *UserService) SignIn(ctx context.Context, username string, password string) (*User, error) {
	la := LoginAttempt{
		Username: strings.ToLower(username),
		Password: password,
	}
	// err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return s.userRepo.signIn(ctx, la)
}

func (s *UserService) SignUp(ctx context.Context, username string, password string, name string) error {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		return err
	}

	sa := SignupAttempt{
		Username: strings.ToLower(username),
		Password: string(passwordBytes),
		Name:     name,
	}
	return s.userRepo.signUp(ctx, sa)
}
