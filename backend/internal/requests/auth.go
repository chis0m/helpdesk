package requests

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,max=120"`
	Password string `json:"password" binding:"required,min=8,max=128"`
}

type SignupRequest struct {
	Email      string  `json:"email" binding:"required,email,max=120"`
	Password   string  `json:"password" binding:"required,min=8,max=128"`
	FirstName  string  `json:"first_name" binding:"required,min=2,max=100"`
	LastName   string  `json:"last_name" binding:"required,min=2,max=100"`
	MiddleName *string `json:"middle_name" binding:"omitempty,max=100"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required,min=8,max=128"`
	NewPassword     string `json:"new_password" binding:"required,min=8,max=128"`
}
