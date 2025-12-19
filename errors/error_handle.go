package errors

import "errors"

// Common errors used across service & repository

var (
	ErrBadRequest   = errors.New("bad request")
	ErrNotFound     = errors.New("resource not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrInternal     = errors.New("internal error")

	// Ticket-specific
	ErrTicketNotFound = errors.New("ticket not found")
	ErrInvalidInput   = errors.New("invalid input")
)
