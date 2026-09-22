package repository

import (
	"context"

	"7solutions-challenge/user/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository interface {
	Create(ctx context.Context, user model.User) error
	GetByID(ctx context.Context, id string) (model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
	List(ctx context.Context) ([]model.User, error)
	Update(ctx context.Context, user model.User) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
}

type userRepository struct {
	col *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	collection := db.Collection("users")

	repo := &userRepository{col: collection}
	err := repo.createUserIndexes()

	if err != nil {
		panic(err)
	}

	return repo
}

func (r *userRepository) createUserIndexes() error {
	_, err := r.col.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.D{
				{Key: "email", Value: 1},
			},
			Options: options.Index().
				SetUnique(true),
		},
	)

	return err
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

func (r *userRepository) List(ctx context.Context) (res []model.User, err error) {
	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	if err := cur.All(ctx, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (r *userRepository) Update(ctx context.Context, user model.User) error {
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{
		"$set": bson.M{
			"name":  user.Name,
			"email": user.Email,
		},
	})

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return err
		}
		return err
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	res, err := r.col.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return err
	}
	return nil
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	return r.col.CountDocuments(ctx, bson.M{})
}
