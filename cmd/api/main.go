package main

import (
	"7solutions-challenge/config"
	"7solutions-challenge/database"
	"7solutions-challenge/helper/hash"
	"7solutions-challenge/helper/jwt"
	"7solutions-challenge/router"
	"7solutions-challenge/user"
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	jwtManager, err := jwt.NewJwtManager(cfg.JwtSecret, cfg.JwtTtl)
	if err != nil {
		panic(err)
	}

	bcrypt := hash.NewBcrypt()
	mongo, err := database.ConnectMongo(cfg.MongoURI)
	db := mongo.Database(cfg.MongoDB)

	r := gin.New()
	router.Register(r, cfg, jwtManager, db, bcrypt)

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM,
	)
	defer cancel()

	userWorker := user.NewWorker(db, bcrypt)
	if err != nil {
		log.Fatal(err)
	}

	go userWorker.Run(ctx)

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	r.Run()
}
