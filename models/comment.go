package models

import "time"

type NewComment struct {
	ID          int64     `json:"id"`
	UUID        string    `json:"uuid"`

	TicketID    int64     `json:"ticket_id"`
	CreatedByID int64     `json:"created_by_id"`

	CommentText string    `json:"comment_text"`
	CreatedAt   time.Time `json:"created_at"`
}