package model

type LoginResponse struct {
	Token     string   `json:"token"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Roles     []string `json:"roles"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
