package service

import (
	"7solutions-challenge/auth/domain"
	"7solutions-challenge/auth/model"
	"context"
	"errors"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockUserRepository struct {
	createFn     func(context.Context, model.User) error
	getByEmailFn func(context.Context, string) (model.User, error)
	createdUser  model.User
}

func (m *mockUserRepository) Create(
	ctx context.Context,
	user model.User,
) error {
	m.createdUser = user

	if m.createFn != nil {
		return m.createFn(ctx, user)
	}

	return nil
}

func (m *mockUserRepository) GetByID(
	ctx context.Context,
	id string,
) (model.User, error) {
	return model.User{}, nil
}

func (m *mockUserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (model.User, error) {
	if m.getByEmailFn != nil {
		return m.getByEmailFn(ctx, email)
	}

	return model.User{}, nil
}

func (m *mockUserRepository) List(
	ctx context.Context,
) ([]model.User, error) {
	return nil, nil
}

func (m *mockUserRepository) Update(
	ctx context.Context,
	user model.User,
) error {
	return nil
}

func (m *mockUserRepository) Delete(
	ctx context.Context,
	id string,
) error {
	return nil
}

func (m *mockUserRepository) Count(
	ctx context.Context,
) (int64, error) {
	return 0, nil
}

type mockBcryptHasher struct {
	hashFn    func(string) (string, error)
	compareFn func(string, string) error
}

func (m *mockBcryptHasher) Hash(password string) (string, error) {
	if m.hashFn != nil {
		return m.hashFn(password)
	}

	return "hashed-password", nil
}

func (m *mockBcryptHasher) Compare(
	hashedPassword string,
	password string,
) error {
	if m.compareFn != nil {
		return m.compareFn(hashedPassword, password)
	}

	return nil
}

type mockJWTManager struct {
	generateFn func(string) (string, error)
}

func (m *mockJWTManager) Generate(userID string) (string, error) {
	if m.generateFn != nil {
		return m.generateFn(userID)
	}

	return "test-token", nil
}

func (m *mockJWTManager) Parse(tokenString string) (string, error) {
	return "", nil
}

func TestAuthService_Register(t *testing.T) {
	repo := &mockUserRepository{}

	bcrypt := &mockBcryptHasher{
		hashFn: func(password string) (string, error) {
			if password != "password123" {
				t.Errorf("unexpected password: %s", password)
			}

			return "hashed-password", nil
		},
	}

	jwtManager := &mockJWTManager{}

	service := NewAuthService(
		repo,
		jwtManager,
		bcrypt,
	)

	req := domain.UserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	err := service.Register(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.createdUser.Name != req.Name {
		t.Errorf(
			"expected name %q, got %q",
			req.Name,
			repo.createdUser.Name,
		)
	}

	if repo.createdUser.Email != req.Email {
		t.Errorf(
			"expected email %q, got %q",
			req.Email,
			repo.createdUser.Email,
		)
	}

	if repo.createdUser.Password != "hashed-password" {
		t.Errorf(
			"expected hashed password, got %q",
			repo.createdUser.Password,
		)
	}

	if repo.createdUser.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestAuthService_Register_HashError(t *testing.T) {
	expectedErr := errors.New("hash error")

	repo := &mockUserRepository{}

	bcrypt := &mockBcryptHasher{
		hashFn: func(password string) (string, error) {
			return "", expectedErr
		},
	}

	service := NewAuthService(
		repo,
		&mockJWTManager{},
		bcrypt,
	)

	err := service.Register(
		context.Background(),
		domain.UserRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		},
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf(
			"expected %v, got %v",
			expectedErr,
			err,
		)
	}
}

func TestAuthService_Register_RepositoryError(t *testing.T) {
	expectedErr := errors.New("repository error")

	repo := &mockUserRepository{
		createFn: func(
			ctx context.Context,
			user model.User,
		) error {
			return expectedErr
		},
	}

	service := NewAuthService(
		repo,
		&mockJWTManager{},
		&mockBcryptHasher{},
	)

	err := service.Register(
		context.Background(),
		domain.UserRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "password123",
		},
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf(
			"expected %v, got %v",
			expectedErr,
			err,
		)
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	userID := primitive.NewObjectID()

	repo := &mockUserRepository{
		getByEmailFn: func(
			ctx context.Context,
			email string,
		) (model.User, error) {
			if email != "john@example.com" {
				t.Errorf(
					"expected email %q, got %q",
					"john@example.com",
					email,
				)
			}

			return model.User{
				ID:       userID,
				Email:    "john@example.com",
				Password: "hashed-password",
			}, nil
		},
	}

	bcrypt := &mockBcryptHasher{
		compareFn: func(
			hashedPassword string,
			password string,
		) error {
			if hashedPassword != "hashed-password" {
				t.Errorf(
					"unexpected hashed password: %s",
					hashedPassword,
				)
			}

			if password != "password123" {
				t.Errorf(
					"unexpected password: %s",
					password,
				)
			}

			return nil
		},
	}

	jwtManager := &mockJWTManager{
		generateFn: func(userIDArg string) (string, error) {
			if userIDArg != userID.Hex() {
				t.Errorf(
					"expected user ID %q, got %q",
					userID.Hex(),
					userIDArg,
				)
			}

			return "jwt-token", nil
		},
	}

	service := NewAuthService(
		repo,
		jwtManager,
		bcrypt,
	)

	token, err := service.Login(
		context.Background(),
		domain.LoginRequest{
			Email:    "john@example.com",
			Password: "password123",
		},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token != "jwt-token" {
		t.Errorf(
			"expected token %q, got %q",
			"jwt-token",
			token,
		)
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	repoErr := errors.New("user not found")

	repo := &mockUserRepository{
		getByEmailFn: func(
			ctx context.Context,
			email string,
		) (model.User, error) {
			return model.User{}, repoErr
		},
	}

	service := NewAuthService(
		repo,
		&mockJWTManager{},
		&mockBcryptHasher{},
	)

	_, err := service.Login(
		context.Background(),
		domain.LoginRequest{
			Email:    "unknown@example.com",
			Password: "password123",
		},
	)

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf(
			"expected %v, got %v",
			ErrInvalidCredentials,
			err,
		)
	}
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	repo := &mockUserRepository{
		getByEmailFn: func(
			ctx context.Context,
			email string,
		) (model.User, error) {
			return model.User{
				ID:       primitive.NewObjectID(),
				Email:    email,
				Password: "hashed-password",
			}, nil
		},
	}

	bcrypt := &mockBcryptHasher{
		compareFn: func(
			hashedPassword string,
			password string,
		) error {
			return errors.New("password mismatch")
		},
	}

	service := NewAuthService(
		repo,
		&mockJWTManager{},
		bcrypt,
	)

	_, err := service.Login(
		context.Background(),
		domain.LoginRequest{
			Email:    "john@example.com",
			Password: "wrong-password",
		},
	)

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf(
			"expected %v, got %v",
			ErrInvalidCredentials,
			err,
		)
	}
}

func TestAuthService_Login_JWTError(t *testing.T) {
	expectedErr := errors.New("jwt generation failed")
	userID := primitive.NewObjectID()

	repo := &mockUserRepository{
		getByEmailFn: func(
			ctx context.Context,
			email string,
		) (model.User, error) {
			return model.User{
				ID:       userID,
				Email:    email,
				Password: "hashed-password",
			}, nil
		},
	}

	jwtManager := &mockJWTManager{
		generateFn: func(userIDArg string) (string, error) {
			return "", expectedErr
		},
	}

	service := NewAuthService(
		repo,
		jwtManager,
		&mockBcryptHasher{},
	)

	_, err := service.Login(
		context.Background(),
		domain.LoginRequest{
			Email:    "john@example.com",
			Password: "password123",
		},
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf(
			"expected %v, got %v",
			expectedErr,
			err,
		)
	}
}
