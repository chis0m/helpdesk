package requests

type ListAuditLogsQuery struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}
