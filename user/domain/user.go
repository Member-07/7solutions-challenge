package domain

import "time"

type UserRequest struct {
	Name     string `json:"name" validate:"required,max=300"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,max=72"`
}

type UserUpdateRequest struct {
	ID    string `json:"id" validate:"required"`
	Name  string `json:"name" validate:"required,max=300"`
	Email string `json:"email" validate:"required,email"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}
