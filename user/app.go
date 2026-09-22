package user

import (
	"7solutions-challenge/helper/hash"
	"7solutions-challenge/user/handler"
	"7solutions-challenge/user/repository"
	"7solutions-challenge/user/service"

	"go.mongodb.org/mongo-driver/mongo"
)

func New(db *mongo.Database, bcrypt hash.BcryptHasher) handler.UserHandler {
	repo := repository.NewUserRepository(db)
	service := service.NewUserService(repo, bcrypt)
	h := handler.NewUserHandler(service)

	return h
}
