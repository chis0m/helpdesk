package requests

type CreateStaffInviteRequest struct {
	Email      string  `json:"email" binding:"required,email,max=120"`
	FirstName  string  `json:"first_name" binding:"required,min=2,max=100"`
	LastName   string  `json:"last_name" binding:"required,min=2,max=100"`
	MiddleName *string `json:"middle_name" binding:"omitempty,max=100"`
	// Role is the account role when the invite is accepted: staff (default) or admin. Only admin or super_admin may specify admin.
	Role string `json:"role" binding:"omitempty,oneof=staff admin"`
}

type AcceptInviteRequest struct {
	Token    string `json:"token" binding:"required,len=64,hexadecimal"`
	Password string `json:"password" binding:"required,min=8,max=128"`
}
