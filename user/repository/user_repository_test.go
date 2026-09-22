package repository

import (
	"7solutions-challenge/user/model"
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

func TestUserRepository_List(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		id1 := primitive.NewObjectID()
		id2 := primitive.NewObjectID()

		mt.AddMockResponses(
			mtest.CreateCursorResponse(
				1,
				"users.users",
				mtest.FirstBatch,
				bson.D{
					{Key: "_id", Value: id1},
					{Key: "name", Value: "John"},
					{Key: "email", Value: "john@example.com"},
				},
				bson.D{
					{Key: "_id", Value: id2},
					{Key: "name", Value: "Jane"},
					{Key: "email", Value: "jane@example.com"},
				},
			),
			mtest.CreateCursorResponse(
				0,
				"users.users",
				mtest.NextBatch,
			),
		)

		repo := &userRepository{col: mt.Coll}

		users, err := repo.List(context.Background())
		if err != nil {
			mt.Fatalf("expected no error, got %v", err)
		}

		if len(users) != 2 {
			mt.Fatalf("expected 2 users, got %d", len(users))
		}

		if users[0].Name != "John" {
			t.Errorf("expected John, got %s", users[0].Name)
		}

		if users[1].Name != "Jane" {
			t.Errorf("expected Jane, got %s", users[1].Name)
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

		_, err := repo.List(context.Background())
		if err == nil {
			mt.Fatal("expected error, got nil")
		}
	})
}

func TestUserRepository_Update(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(
			bson.D{
				{Key: "ok", Value: 1},
				{Key: "n", Value: int64(1)},
				{Key: "nModified", Value: int64(1)},
			},
		)

		repo := &userRepository{col: mt.Coll}

		user := model.User{
			ID:    primitive.NewObjectID(),
			Name:  "John Updated",
			Email: "john.updated@example.com",
		}

		err := repo.Update(context.Background(), user)
		if err != nil {
			mt.Fatalf("expected no error, got %v", err)
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

		err := repo.Update(context.Background(), model.User{
			ID: primitive.NewObjectID(),
		})

		if err == nil {
			mt.Fatal("expected error, got nil")
		}
	})
}

func TestUserRepository_Delete(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(
			bson.D{
				{Key: "ok", Value: 1},
				{Key: "n", Value: int64(1)},
			},
		)

		repo := &userRepository{col: mt.Coll}

		err := repo.Delete(
			context.Background(),
			primitive.NewObjectID().Hex(),
		)

		if err != nil {
			mt.Fatalf("expected no error, got %v", err)
		}
	})

	mt.Run("invalid id", func(mt *mtest.T) {
		repo := &userRepository{col: mt.Coll}

		err := repo.Delete(
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

		err := repo.Delete(
			context.Background(),
			primitive.NewObjectID().Hex(),
		)

		if err == nil {
			mt.Fatal("expected error, got nil")
		}
	})

}

func TestUserRepository_Count(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(
			mtest.CreateCursorResponse(
				1,
				"users.users",
				mtest.FirstBatch,
				bson.D{
					{Key: "n", Value: int64(5)},
				},
			),
		)

		repo := &userRepository{col: mt.Coll}

		count, err := repo.Count(context.Background())
		if err != nil {
			mt.Fatalf("expected no error, got %v", err)
		}

		if count != 5 {
			mt.Errorf("expected count 5, got %d", count)
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

		_, err := repo.Count(context.Background())
		if err == nil {
			mt.Fatal("expected error, got nil")
		}
	})
}
