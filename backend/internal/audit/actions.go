package audit

// Stable action names for queries and incident response.
const (
	ActionAPIMutation = "api.mutation"

	ActionCSRFDenied = "security.csrf_denied"

	ActionAuthLoginSuccess = "auth.login.success"
	ActionAuthLoginFailure = "auth.login.failure"

	ActionAuthLogout = "auth.logout"

	ActionAuthPasswordChangeSuccess = "auth.password_change.success"
	ActionAuthPasswordChangeFailure = "auth.password_change.failure"

	ActionAuthSessionRevoke        = "auth.session.revoke"
	ActionAuthSessionsRevokeOthers = "auth.sessions.revoke_others"

	ActionAdminRoleUpdate   = "admin.user.role_update"
	ActionAdminStaffCreate  = "admin.staff.create"
	ActionAdminAuditLogList = "admin.audit_log.list"
	ActionInviteStaffCreate = "admin.invite.staff_create"
	ActionInviteAccepted    = "invite.accept"

	ActionTicketCreate       = "ticket.create"
	ActionTicketSearch       = "ticket.search"
	ActionTicketUpdate       = "ticket.update"
	ActionTicketDelete       = "ticket.delete"
	ActionTicketStatusUpdate = "ticket.status_update"
	ActionTicketAssign       = "ticket.assign"

	ActionTicketCommentCreate = "ticket_comment.create"
	ActionTicketCommentUpdate = "ticket_comment.update"
	ActionTicketCommentDelete = "ticket_comment.delete"

	ActionUserProfileUpdate = "user.profile.update"
	ActionUserDelete        = "user.delete"
)

// ResourceType values for ResourceType column (stable identifiers).
const (
	ResourceTypeTicket        = "ticket"
	ResourceTypeUser          = "user"
	ResourceTypeTicketComment = "ticket_comment"
	ResourceTypeInvite        = "invite"
	ResourceTypeAuditLog      = "audit_log"
)
