package service

import (
	"7solutions-challenge/helper/hash"
	"7solutions-challenge/user/domain"
	"7solutions-challenge/user/model"
	"context"
	"errors"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockUserRepository struct {
	createFn  func(context.Context, model.User) error
	getByIDFn func(context.Context, string) (model.User, error)
	listFn    func(context.Context) ([]model.User, error)
	updateFn  func(context.Context, model.User) error
	deleteFn  func(context.Context, string) error
	countFn   func(context.Context) (int64, error)

	createdUser model.User
	updatedUser model.User
	deletedID   string
}

func (m *mockUserRepository) Create(ctx context.Context, user model.User) error {
	m.createdUser = user
	if m.createFn != nil {
		return m.createFn(ctx, user)
	}
	return nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, id string) (model.User, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return model.User{}, nil
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (model.User, error) {
	return model.User{}, nil
}

func (m *mockUserRepository) List(ctx context.Context) ([]model.User, error) {
	if m.listFn != nil {
		return m.listFn(ctx)
	}
	return nil, nil
}

func (m *mockUserRepository) Update(ctx context.Context, user model.User) error {
	m.updatedUser = user
	if m.updateFn != nil {
		return m.updateFn(ctx, user)
	}
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id string) error {
	m.deletedID = id
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func (m *mockUserRepository) Count(ctx context.Context) (int64, error) {
	if m.countFn != nil {
		return m.countFn(ctx)
	}
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

func (m *mockBcryptHasher) Compare(hashedPassword, password string) error {
	if m.compareFn != nil {
		return m.compareFn(hashedPassword, password)
	}
	return nil
}

var _ hash.BcryptHasher = (*mockBcryptHasher)(nil)

func TestUserService_Create(t *testing.T) {
	repo := &mockUserRepository{}

	bcrypt := &mockBcryptHasher{
		hashFn: func(password string) (string, error) {
			return "hashed-password", nil
		},
	}

	service := NewUserService(repo, bcrypt)

	req := domain.UserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	err := service.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.createdUser.Name != req.Name {
		t.Errorf("expected name %q, got %q", req.Name, repo.createdUser.Name)
	}

	if repo.createdUser.Email != req.Email {
		t.Errorf("expected email %q, got %q", req.Email, repo.createdUser.Email)
	}

	if repo.createdUser.Password != "hashed-password" {
		t.Errorf("expected hashed password, got %q", repo.createdUser.Password)
	}

	if repo.createdUser.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestUserService_Create_HashError(t *testing.T) {
	expectedErr := errors.New("hash failed")

	repo := &mockUserRepository{}

	bcrypt := &mockBcryptHasher{
		hashFn: func(password string) (string, error) {
			return "", expectedErr
		},
	}

	service := NewUserService(repo, bcrypt)

	err := service.Create(context.Background(), domain.UserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}

func TestUserService_Create_RepositoryError(t *testing.T) {
	expectedErr := errors.New("repository error")

	repo := &mockUserRepository{
		createFn: func(ctx context.Context, user model.User) error {
			return expectedErr
		},
	}

	service := NewUserService(repo, &mockBcryptHasher{})

	err := service.Create(context.Background(), domain.UserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}

func TestUserService_GetByID(t *testing.T) {
	userID := "507f1f77bcf86cd799439011"

	expectedUser := model.User{
		ID:        mustObjectID(userID),
		Name:      "John Doe",
		Email:     "john@example.com",
		Password:  "hashed-password",
		CreatedAt: time.Now(),
	}

	repo := &mockUserRepository{
		getByIDFn: func(ctx context.Context, id string) (model.User, error) {
			return expectedUser, nil
		},
	}

	service := NewUserService(repo, &mockBcryptHasher{})

	res, err := service.GetByID(context.Background(), userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.ID != userID {
		t.Errorf("expected ID %q, got %q", userID, res.ID)
	}

	if res.Name != expectedUser.Name {
		t.Errorf("expected name %q, got %q", expectedUser.Name, res.Name)
	}

	if res.Email != expectedUser.Email {
		t.Errorf("expected email %q, got %q", expectedUser.Email, res.Email)
	}

	if res.Password != "" {
		t.Error("password should not be exposed")
	}
}

func TestUserService_GetByID_Error(t *testing.T) {
	expectedErr := errors.New("user not found")

	repo := &mockUserRepository{
		getByIDFn: func(ctx context.Context, id string) (model.User, error) {
			return model.User{}, expectedErr
		},
	}

	service := NewUserService(repo, &mockBcryptHasher{})

	_, err := service.GetByID(
		context.Background(),
		"507f1f77bcf86cd799439011",
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}

func TestUserService_List(t *testing.T) {
	users := []model.User{
		{
			ID:    mustObjectID("507f1f77bcf86cd799439011"),
			Name:  "John Doe",
			Email: "john@example.com",
		},
		{
			ID:    mustObjectID("507f1f77bcf86cd799439012"),
			Name:  "Jane Doe",
			Email: "jane@example.com",
		},
	}

	repo := &mockUserRepository{
		listFn: func(ctx context.Context) ([]model.User, error) {
			return users, nil
		},
	}

	service := NewUserService(repo, &mockBcryptHasher{})

	res, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(res) != 2 {
		t.Fatalf("expected 2 users, got %d", len(res))
	}

	if res[0].Name != "John Doe" {
		t.Errorf("unexpected first user: %+v", res[0])
	}

	if res[1].Email != "jane@example.com" {
		t.Errorf("unexpected second user: %+v", res[1])
	}
}

func TestUserService_Update(t *testing.T) {
	userID := "507f1f77bcf86cd799439011"

	existingUser := model.User{
		ID:        mustObjectID(userID),
		Name:      "Old Name",
		Email:     "old@example.com",
		Password:  "hashed-password",
		CreatedAt: time.Now(),
	}

	repo := &mockUserRepository{
		getByIDFn: func(ctx context.Context, id string) (model.User, error) {
			return existingUser, nil
		},
	}

	service := NewUserService(repo, &mockBcryptHasher{})

	req := domain.UserUpdateRequest{
		ID:    userID,
		Name:  "New Name",
		Email: "new@example.com",
	}

	err := service.Update(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.updatedUser.Name != "New Name" {
		t.Errorf("expected name %q, got %q", "New Name", repo.updatedUser.Name)
	}

	if repo.updatedUser.Email != "new@example.com" {
		t.Errorf("expected email %q, got %q", "new@example.com", repo.updatedUser.Email)
	}

	if repo.updatedUser.Password != existingUser.Password {
		t.Error("password should not be changed")
	}
}

func TestUserService_Delete(t *testing.T) {
	repo := &mockUserRepository{}

	service := NewUserService(repo, &mockBcryptHasher{})

	userID := "507f1f77bcf86cd799439011"

	err := service.Delete(context.Background(), userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.deletedID != userID {
		t.Errorf("expected deleted ID %q, got %q", userID, repo.deletedID)
	}
}

func TestUserService_Count(t *testing.T) {
	repo := &mockUserRepository{
		countFn: func(ctx context.Context) (int64, error) {
			return 10, nil
		},
	}

	service := NewUserService(repo, &mockBcryptHasher{})

	count, err := service.Count(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if count != 10 {
		t.Errorf("expected count 10, got %d", count)
	}
}

func mustObjectID(id string) primitive.ObjectID {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}

	return objectID
}
