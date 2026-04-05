package repositories

import (
	"errors"

	"gorm.io/gorm"

	"helpdesk/backend/internal/models"
	"helpdesk/backend/internal/requests"
)

type TicketCommentRepository struct {
	db *gorm.DB
}

func NewTicketCommentRepository(db *gorm.DB) *TicketCommentRepository {
	return &TicketCommentRepository{db: db}
}

func (r *TicketCommentRepository) Create(input requests.CreateTicketCommentInput) (*models.TicketComment, error) {
	comment := &models.TicketComment{
		TicketID:     input.TicketID,
		AuthorUserID: input.AuthorUserID,
		Body:         input.Body,
	}
	if err := r.db.Create(comment).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

func (r *TicketCommentRepository) ListByTicketID(ticketID uint64) ([]models.TicketCommentWithAuthor, error) {
	var comments []models.TicketCommentWithAuthor
	err := r.db.Table("ticket_comments AS tc").
		Select(`
			tc.id AS comment_id,
			tc.ticket_id,
			tc.author_user_id,
			u.email AS author_email,
			u.first_name AS author_first_name,
			u.last_name AS author_last_name,
			tc.body,
			tc.created_at,
			tc.updated_at
		`).
		Joins("JOIN users AS u ON u.id = tc.author_user_id").
		Where("tc.ticket_id = ?", ticketID).
		Order("tc.created_at ASC").
		Scan(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *TicketCommentRepository) GetByIDAndTicketID(commentID uint64, ticketID uint64) (*models.TicketComment, error) {
	var comment models.TicketComment
	if err := r.db.First(&comment, "id = ? AND ticket_id = ?", commentID, ticketID).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *TicketCommentRepository) UpdateBody(commentID uint64, body string) (*models.TicketComment, error) {
	var comment models.TicketComment
	if err := r.db.First(&comment, "id = ?", commentID).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&comment).Update("body", body).Error; err != nil {
		return nil, err
	}

	if err := r.db.First(&comment, "id = ?", commentID).Error; err != nil {
		return nil, errors.New("comment updated but failed to reload")
	}
	return &comment, nil
}

func (r *TicketCommentRepository) DeleteByID(commentID uint64) error {
	return r.db.Delete(&models.TicketComment{}, "id = ?", commentID).Error
}
