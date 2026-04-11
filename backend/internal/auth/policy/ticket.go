package policy

import "helpdesk/backend/internal/models"

// SEC-02: Central ticket authorization — admin/super_admin, or reporter/assignee only.
func CanAccessTicket(actor *models.User, actorRole models.UserRole, ticket *models.Ticket) bool {
	if actor == nil || ticket == nil {
		return false
	}
	switch actorRole {
	case models.RoleAdmin, models.RoleSuperAdmin:
		return true
	default:
		if ticket.ReporterUserID == actor.ID {
			return true
		}
		if ticket.AssignedUserID != nil && *ticket.AssignedUserID == actor.ID {
			return true
		}
		return false
	}
}
