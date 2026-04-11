package audit

// Event describes one append-only audit row. Identity comes from the server auth context.
type Event struct {
	Action       string
	Success      bool
	ErrorCode    string
	HTTPMethod   string
	Path         string
	ActorUserID  *uint64
	SessionID    *string
	TokenJTI     *string
	IP           *string
	UserAgent    *string
	ResourceType *string
	ResourceID   *uint64
	Metadata     map[string]interface{}
}
