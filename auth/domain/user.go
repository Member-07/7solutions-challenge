package domain

type UserRequest struct {
	Name     string `json:"name" validate:"required,max=300"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,max=72"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,max=72"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
