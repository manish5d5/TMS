package repos

import (
	"context"

	"TMS/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentRepo struct {
	db *pgxpool.Pool
}

func NewCommentRepo(db *pgxpool.Pool) *CommentRepo {
	return &CommentRepo{db: db}
}



func (r *CommentRepo) CreateComment(
	ctx context.Context,
	c models.NewComment,
) (int64, error) {

	query := `
		INSERT INTO comments (ticket_id, created_by_id, comment_text)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var id int64
	err := r.db.QueryRow(
		ctx,
		query,
		c.TicketID,
		c.CreatedByID,
		c.CommentText,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}



func (r *CommentRepo) GetCommentsByTicket(
	ctx context.Context,
	ticketID int64,
) ([]models.NewComment, error) {

	query := `
		SELECT id, uuid, ticket_id, created_by_id, comment_text, created_at
		FROM comments
		WHERE ticket_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []models.NewComment{}

	for rows.Next() {
		var c models.NewComment
		err := rows.Scan(
			&c.ID,
			&c.UUID,
			&c.TicketID,
			&c.CreatedByID,
			&c.CommentText,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}


