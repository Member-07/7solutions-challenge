package repository

import (
	"context"

	"7solutions-challenge/auth/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(ctx context.Context, user model.User) error
	GetByID(ctx context.Context, id string) (model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
}

type userRepository struct {
	col *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	collection := db.Collection("users")

	return &userRepository{col: collection}
}

func (r *userRepository) Create(ctx context.Context, user model.User) error {
	if _, err := r.col.InsertOne(ctx, user); err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (res model.User, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return res, err
	}

	err = r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (res model.User, err error) {
	err = r.col.FindOne(ctx, bson.M{"email": email}).Decode(&res)

	if err != nil {
		return res, err
	}
	return res, nil
}
