package requests

type CreateTicketCommentRequest struct {
	// SEC-03: After bind, service applies strings.TrimSpace on body; angle brackets preserved. UI escapes output.
	Body string `json:"body" binding:"required,min=1,max=5000"`
}

type CreateTicketCommentInput struct {
	TicketID     uint64
	AuthorUserID uint64
	Body         string
}

type UpdateTicketCommentRequest struct {
	// SEC-03: After bind, service applies strings.TrimSpace on body; angle brackets preserved. UI escapes output.
	Body string `json:"body" binding:"required,min=1,max=5000"`
}
