package services

import (
	"context"
	"errors"

	"TMS/models"
	"TMS/repos"
)

type TicketsService struct {
	repo repos.TicketsRepo
}

func NewTicketsService(r repos.TicketsRepo) *TicketsService {
	return &TicketsService{repo: r}
}

func (s *TicketsService) CreateTicket(ctx context.Context,input models.NewTicketFormat) (int, error) {

	if input.Title == "" || input.Description == "" {
		return 0, errors.New("title and description are required")
	}

	if input.Priority == "" {
		input.Priority = "medium"
	}

	userId:=input.CreatedBy
	if userId==0{
		return 0,errors.New("user not exsist")
	}//need to be modify

	ticket := models.Ticket{
		Title:       input.Title,
		Description: input.Description,
		Priority:    input.Priority,
		Status:      "open",
		CreatedBy:   userId,
	}

	return s.repo.CreateTicket(ctx, ticket)
}



func (s *TicketsService) GetTicketById(ctx context.Context,id int) (*models.Ticket, error) {
	ticket, err := s.repo.GetTicketById(ctx, id)
	if err != nil {
		return nil, err
	}

	if ticket == nil {
		return nil,  errors.New("ticket not found")
	}

	return ticket, nil
}



func (s *TicketsService) DeleteTicketById(ctx context.Context,ticketID int) error {

	// userID, ok := ctx.Value("user_id").(int)
	// if !ok {
	// 	return  errors.New("invalid user")
	// }

	// ownerID, err := s.repo.GetTicketOwner(ctx, ticketID)
	// if err != nil {
	// 	return err
	// }

	// if ownerID != userID {
	// 	return apperror.Forbidden("you are not allowed to delete this ticket")
	// }

	return s.repo.DeleteTicketById(ctx, ticketID)
}



func (s *TicketsService) GetTicketsByFilters(
	ctx context.Context,
	status string,
	priority string,
) ([]models.Ticket, error) {

	// Optional validation
	if status != "" && status != "open" && status != "closed" {
		return nil,  errors.New("invalid status value")
	}

	if priority != "" &&
		priority != "low" &&
		priority != "medium" &&
		priority != "high" {
		return nil,  errors.New("invalid priority value")
	}

	return s.repo.GetTicketsByFilters(ctx, status, priority)
}



func (s *TicketsService) UpdateTicket(
	ctx context.Context,
	id int,
	input models.NewTicketFormat,
) error {

	if input.Title == "" || input.Description == "" {
		return errors.New("title and description are required")
	}

	if input.Priority == "" {
		input.Priority = "medium"
	}

	return s.repo.UpdateTicket(ctx, id, input)
}
