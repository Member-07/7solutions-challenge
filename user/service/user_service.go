package service

import (
	"7solutions-challenge/helper/hash"
	"7solutions-challenge/user/domain"
	"7solutions-challenge/user/model"
	"7solutions-challenge/user/repository"
	"context"
	"time"
)

type UserService interface {
	Create(ctx context.Context, req domain.UserRequest) error
	GetByID(ctx context.Context, id string) (domain.UserResponse, error)
	List(ctx context.Context) ([]domain.UserResponse, error)
	Update(ctx context.Context, req domain.UserUpdateRequest) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
}

type userService struct {
	repo   repository.UserRepository
	bcrypt hash.BcryptHasher
}

func NewUserService(repo repository.UserRepository, bcrypt hash.BcryptHasher) UserService {
	return &userService{
		repo:   repo,
		bcrypt: bcrypt,
	}
}
func (s *userService) Create(ctx context.Context, req domain.UserRequest) error {

	hashedPassword, err := s.bcrypt.Hash(req.Password)
	if err != nil {
		return err
	}

	user := model.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}
	return s.repo.Create(ctx, user)
}

func (s *userService) GetByID(ctx context.Context, id string) (domain.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.UserResponse{}, err
	}
	return toUserResponse(user), nil
}

func (s *userService) List(ctx context.Context) ([]domain.UserResponse, error) {
	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]domain.UserResponse, 0, len(users))
	for _, user := range users {
		res = append(res, toUserResponse(user))
	}
	return res, nil
}

func (s *userService) Update(ctx context.Context, req domain.UserUpdateRequest) error {
	user, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}

	user.Name = req.Name
	user.Email = req.Email
	return s.repo.Update(ctx, user)
}

func (s *userService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *userService) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

func toUserResponse(user model.User) domain.UserResponse {
	return domain.UserResponse{
		ID:        user.ID.Hex(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}
