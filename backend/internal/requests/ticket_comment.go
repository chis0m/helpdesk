package requests

type CreateTicketCommentRequest struct {
	// VULN-04: Input validation is weak; stored XSS is possible.
	Body string `json:"body" binding:"required,min=1,max=5000"`
}

type CreateTicketCommentInput struct {
	TicketID     uint64
	AuthorUserID uint64
	Body         string
}

type UpdateTicketCommentRequest struct {
	// VULN-04: Input validation is weak; stored XSS is possible.
	Body string `json:"body" binding:"required,min=1,max=5000"`
}
