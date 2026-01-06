package models

type RegisterRequest struct {
	UserName        string `json:"userName"`
	Email           string `json:"email"`
	NoHandphone     string `json:"noHandphone"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type AuthResponse struct {
	UserName string `json:"username"`
	Message  string `json:"message"`
}

type UserInfo struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
