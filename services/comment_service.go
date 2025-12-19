package services

import (
	"context"
	"errors"

	"TMS/models"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, c models.NewComment) (int64, error)
	GetCommentsByTicket(ctx context.Context, ticketID int64) ([]models.NewComment, error)
}

type CommentService struct {
	repo CommentRepository
}

func NewCommentService(r CommentRepository) *CommentService {
	return &CommentService{repo: r}
}


func (s *CommentService) CreateComment(
	ctx context.Context,
	ticketID int64,
	input models.NewComment,
) (int64, error) {

	if input.CommentText == "" {
		return 0, errors.New("comment text is required")
	}

	// later replace with user_id from JWT context
	userID := int64(1)

	comment := models.NewComment{
		TicketID:    ticketID,
		CreatedByID: userID,
		CommentText: input.CommentText,
	}

	return s.repo.CreateComment(ctx, comment)
}


func (s *CommentService) GetCommentsByTicket(
	ctx context.Context,
	ticketID int64,
) ([]models.NewComment, error) {

	if ticketID <= 0 {
		return nil, errors.New("invalid ticket id")
	}

	return s.repo.GetCommentsByTicket(ctx, ticketID)
}
