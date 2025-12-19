package repos

import (
	"context"
	"errors"
	"fmt"

	"TMS/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TicketsRepo struct {
	db *pgxpool.Pool
}

func NewTicketsRepo(db *pgxpool.Pool) *TicketsRepo {
	return &TicketsRepo{db: db}
}





func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}

func (r *TicketsRepo) CreateTicket(ctx context.Context,t models.Ticket,) (int, error) {

	query := `
		INSERT INTO tickets (title, description, status, priority, created_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var id int
	err := r.db.QueryRow(
		ctx,
		query,
		t.Title,
		t.Description,
		t.Status,
		t.Priority,
		t.CreatedBy,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

	


func (r *TicketsRepo) GetTicketById(ctx context.Context,id int,) (*models.Ticket, error) {

	query := `
		SELECT id, title, description, status, priority,
		       created_by, created_at, updated_at
		FROM tickets
		WHERE id = $1
	`

	var t models.Ticket

	err := r.db.QueryRow(ctx, query, id).Scan(
		&t.ID,
		&t.Title,
		&t.Description,
		&t.Status,
		&t.Priority,
		&t.CreatedBy,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &t, nil
}


func (r *TicketsRepo) DeleteTicketById(ctx context.Context,id int,) error {

	query := `DELETE FROM tickets WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("ticket not found")
	}

	return nil
}




func (r *TicketsRepo) GetTicketsByFilters(
ctx context.Context,
status string,
priority string,
) ([]models.Ticket, error) {

	query := `
		SELECT id, title, description, status, priority,
		       created_by, created_at, updated_at
		FROM tickets
		WHERE 1=1
	`
	args := []any{}
	argPos := 1

	if status != "" {
		query += " AND status = $" + itoa(argPos)
		args = append(args, status)
		argPos++
	}

	if priority != "" {
		query += " AND priority = $" + itoa(argPos)
		args = append(args, priority)
		argPos++
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tickets := []models.Ticket{}

	for rows.Next() {
		var t models.Ticket
		err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.Priority,
			&t.CreatedBy,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}

	return tickets, nil
}


func (r *TicketsRepo) GetTicketOwner(
	ctx context.Context,
	ticketID int,
) (int, error) {

	query := `SELECT created_by FROM tickets WHERE id = $1`

	var ownerID int
	err := r.db.QueryRow(ctx, query, ticketID).Scan(&ownerID)

	if err == pgx.ErrNoRows {
		return 0, errors.New("ticket not found")
	}

	if err != nil {
		return 0, err
	}

	return ownerID, nil
}


func (r *TicketsRepo) UpdateTicket(
	ctx context.Context,
	id int,
	input models.NewTicketFormat,
) error {

	query := `
		UPDATE tickets
		SET
			title = $1,
			description = $2,
			priority = $3,
			updated_at = NOW()
		WHERE id = $4
	`

	cmdTag, err := r.db.Exec(
		ctx,
		query,
		input.Title,
		input.Description,
		input.Priority,
		id,
	)

	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("ticket not found")
	}

	return nil
}
