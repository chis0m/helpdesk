package requests

type CreateTicketCommentRequest struct {
	// VULN-03: Weak input validation / stored XSS risk — no HTML/script sanitization on comment body.
	Body string `json:"body" binding:"required,min=1,max=5000"`
}

type CreateTicketCommentInput struct {
	TicketID     uint64
	AuthorUserID uint64
	Body         string
}

type UpdateTicketCommentRequest struct {
	// VULN-03: Weak input validation / stored XSS risk — no HTML/script sanitization on comment body.
	Body string `json:"body" binding:"required,min=1,max=5000"`
}
