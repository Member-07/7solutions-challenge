package user

import (
	"7solutions-challenge/helper/hash"
	"7solutions-challenge/user/repository"
	"7solutions-challenge/user/service"
	"7solutions-challenge/user/worker"

	"go.mongodb.org/mongo-driver/mongo"
)

func NewWorker(db *mongo.Database, bcrypt hash.BcryptHasher) worker.UserCounter {
	repo := repository.NewUserRepository(db)
	service := service.NewUserService(repo, bcrypt)

	w := worker.NewUserCounter(service)
	return w
}
