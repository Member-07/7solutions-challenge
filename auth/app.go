package auth

import (
	"7solutions-challenge/auth/handler"
	"7solutions-challenge/auth/repository"
	"7solutions-challenge/auth/service"
	"7solutions-challenge/helper/hash"
	"7solutions-challenge/helper/jwt"

	"go.mongodb.org/mongo-driver/mongo"
)

func New(db *mongo.Database, jwt jwt.JWTManager, bcrypt hash.BcryptHasher) handler.AuthHandler {
	repo := repository.NewUserRepository(db)
	service := service.NewAuthService(repo, jwt, bcrypt)
	h := handler.NewAuthHandler(service)

	return h
}
