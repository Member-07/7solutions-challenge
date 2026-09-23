package repository

import (
	"7solutions-challenge/auth/model"
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestUserRepository_Create(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		repo := &userRepository{col: mt.Coll}

		user := model.User{
			ID:        primitive.NewObjectID(),
			Name:      "John Doe",
			Email:     "john@example.com",
			Password:  "hashed-password",
			CreatedAt: time.Now(),
		}

		err := repo.Create(context.Background(), user)
		if err != nil {
			mt.Fatalf("expected no error, got %v", err)
		}
	})

}

func TestUserRepository_GetByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		id := primitive.NewObjectID()

		mt.AddMockResponses(
			mtest.CreateCursorResponse(
				1,
				"users.users",
				mtest.FirstBatch,
				bson.D{
					{Key: "_id", Value: id},
					{Key: "name", Value: "John Doe"},
					{Key: "email", Value: "john@example.com"},
					{Key: "password", Value: "hashed-password"},
				},
			),
		)

		repo := &userRepository{col: mt.Coll}

		user, err := repo.GetByID(
			context.Background(),
			id.Hex(),
		)

		if err != nil {
			mt.Fatalf("expected no error, got %v", err)
		}

		if user.ID != id {
			mt.Errorf("expected ID %v, got %v", id, user.ID)
		}

		if user.Email != "john@example.com" {
			mt.Errorf(
				"expected email %q, got %q",
				"john@example.com",
				user.Email,
			)
		}
	})

	mt.Run("invalid id", func(mt *mtest.T) {
		repo := &userRepository{col: mt.Coll}

		_, err := repo.GetByID(
			context.Background(),
			"invalid-id",
		)

		if err == nil {
			mt.Fatal("expected error, got nil")
		}
	})

	mt.Run("database error", func(mt *mtest.T) {
		mt.AddMockResponses(
			bson.D{
				{Key: "ok", Value: 0},
				{Key: "errmsg", Value: "database error"},
			},
		)

		repo := &userRepository{col: mt.Coll}

		_, err := repo.GetByID(
			context.Background(),
			primitive.NewObjectID().Hex(),
		)

		if err == nil {
			mt.Fatal("expected error, got nil")
		}
	})
}

func TestUserRepository_GetByEmail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		id := primitive.NewObjectID()

		mt.AddMockResponses(
			mtest.CreateCursorResponse(
				1,
				"users.users",
				mtest.FirstBatch,
				bson.D{
					{Key: "_id", Value: id},
					{Key: "name", Value: "John Doe"},
					{Key: "email", Value: "john@example.com"},
					{Key: "password", Value: "hashed-password"},
				},
			),
		)

		repo := &userRepository{col: mt.Coll}

		user, err := repo.GetByEmail(
			context.Background(),
			"john@example.com",
		)

		if err != nil {
			mt.Fatalf("expected no error, got %v", err)
		}

		if user.Email != "john@example.com" {
			mt.Errorf(
				"expected email %q, got %q",
				"john@example.com",
				user.Email,
			)
		}
	})

	mt.Run("not found", func(mt *mtest.T) {
		mt.AddMockResponses(
			mtest.CreateCursorResponse(
				0,
				"users.users",
				mtest.FirstBatch,
			),
		)

		repo := &userRepository{col: mt.Coll}

		_, err := repo.GetByEmail(
			context.Background(),
			"notfound@example.com",
		)

		if err == nil {
			mt.Fatal("expected error, got nil")
		}
	})
}
