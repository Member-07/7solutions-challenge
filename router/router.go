package router

import (
	"7solutions-challenge/auth"
	"7solutions-challenge/config"
	"7solutions-challenge/helper/hash"
	"7solutions-challenge/helper/jwt"
	"7solutions-challenge/middleware"
	"7solutions-challenge/user"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(r *gin.Engine, cfg config.Config, jwt jwt.JWTManager, db *mongo.Database, bcrypt hash.BcryptHasher) {
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	userHandler := user.New(db, bcrypt)
	authHandler := auth.New(db, jwt, bcrypt)

	api := r.Group("/api")
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	user := api.Group("/users")

	user.Use(middleware.Auth(jwt))
	{
		user.POST("/", userHandler.Create)
		user.GET("", userHandler.List)
		user.GET("/:id", userHandler.GetByID)
		user.PUT("/:id", userHandler.Update)
		user.DELETE("/:id", userHandler.Delete)
	}

}
