package https

import (
	"encoding/json"
	"net/http"
	"strconv"

	"TMS/models"
	"TMS/services"

	"github.com/go-chi/chi/v5"
)

type CommentHandler struct {
	tsrv *services.CommentService
}

func NewCommentHandler(s *services.CommentService) *CommentHandler {
	return &CommentHandler{tsrv: s}
}

//
// ─── CREATE NEW COMMENT ───────────────────────────────────────────
func (c *CommentHandler) NewComment(w http.ResponseWriter, r *http.Request) {
	// get ticket id from URL
	ticketIDStr := chi.URLParam(r, "id")

	ticketID, err := strconv.ParseInt(ticketIDStr, 10, 64)
	if err != nil || ticketID <= 0 {
		http.Error(w, "invalid ticket id", http.StatusBadRequest)
		return
	}

	// decode request body
	var data models.NewComment
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// call service
	commentID, err := c.tsrv.CreateComment(r.Context(), ticketID, data)
	if err != nil {
		handleError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, map[string]any{
		"message":    "Comment added successfully",
		"comment_id": commentID,
	})
}



func (c *CommentHandler) GetCommentsByTicket(w http.ResponseWriter, r *http.Request) {
	// get ticket id from URL
	ticketIDStr := chi.URLParam(r, "id")

	ticketID, err := strconv.ParseInt(ticketIDStr, 10, 64)
	if err != nil || ticketID <= 0 {
		http.Error(w, "invalid ticket id", http.StatusBadRequest)
		return
	}

	comments, err := c.tsrv.GetCommentsByTicket(r.Context(), ticketID)
	if err != nil {
		handleError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, comments)
}
