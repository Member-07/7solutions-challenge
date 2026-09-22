package worker

import (
	"7solutions-challenge/user/service"
	"context"
	"log"
	"time"
)

type UserCounter interface {
	Run(ctx context.Context)
}

type userCounter struct {
	service service.UserService
}

func NewUserCounter(service service.UserService) UserCounter {
	return &userCounter{
		service: service,
	}
}

func (w *userCounter) Run(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	log.Println("user counter started")

	for {
		select {
		case <-ticker.C:
			count, err := w.service.Count(ctx)
			if err != nil {
				log.Printf("failed to count users: %v", err)
				continue
			}

			log.Printf("total users: %d", count)

		case <-ctx.Done():
			log.Println("user counter stopped")
			return
		}
	}
}
