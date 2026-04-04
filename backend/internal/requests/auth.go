package requests

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,max=120"`
	Password string `json:"password" binding:"required,min=8,max=128"`
}
