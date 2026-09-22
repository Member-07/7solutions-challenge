package service

import (
	"context"
	"errors"
	"time"

	"7solutions-challenge/auth/domain"
	"7solutions-challenge/auth/model"
	"7solutions-challenge/auth/repository"
	"7solutions-challenge/helper/hash"
	"7solutions-challenge/helper/jwt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type AuthService interface {
	Register(ctx context.Context, req domain.UserRequest) error
	Login(ctx context.Context, req domain.LoginRequest) (string, error)
}

type authService struct {
	repo   repository.UserRepository
	jwt    jwt.JWTManager
	bcrypt hash.BcryptHasher
}

func NewAuthService(
	repo repository.UserRepository,
	jwt jwt.JWTManager,
	bcrypt hash.BcryptHasher,
) AuthService {
	return &authService{
		repo:   repo,
		jwt:    jwt,
		bcrypt: bcrypt,
	}
}

func (s *authService) Register(
	ctx context.Context,
	req domain.UserRequest,
) error {
	password, err := s.bcrypt.Hash(req.Password)
	if err != nil {
		return err
	}

	user := model.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  password,
		CreatedAt: time.Now(),
	}

	return s.repo.Create(ctx, user)
}

func (s *authService) Login(
	ctx context.Context,
	req domain.LoginRequest,
) (string, error) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if err := s.bcrypt.Compare(
		user.Password,
		req.Password,
	); err != nil {
		return "", ErrInvalidCredentials
	}

	return s.jwt.Generate(user.ID.Hex())
}
