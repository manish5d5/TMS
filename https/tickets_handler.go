package https

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"TMS/models"
	"TMS/services"

	"github.com/go-chi/chi/v5"
)

type TicketsHandler struct {
	tsrv *services.TicketsService
}

func NewTicketsHandler(s *services.TicketsService) *TicketsHandler {
	return &TicketsHandler{tsrv: s}
}

//
// ─── HELPER: JSON RESPONSE ───────────────────────────────────────
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

//
// ─── HELPER: ERROR HANDLING ──────────────────────────────────────
func handleError(w http.ResponseWriter, err error) {
	msg := err.Error()

	switch msg {
	case "title and description are required":
		http.Error(w, msg, http.StatusBadRequest)

	case "invalid status value", "invalid priority value":
		http.Error(w, msg, http.StatusBadRequest)

	case "ticket not found":
		http.Error(w, msg, http.StatusNotFound)

	case "user not exist", "invalid user":
		http.Error(w, msg, http.StatusUnauthorized)

	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}


//
// ─── CREATE TICKET ───────────────────────────────────────────────
func (tck *TicketsHandler) NewTicket(w http.ResponseWriter, r *http.Request) {
	var newTck models.NewTicketFormat

	if err := json.NewDecoder(r.Body).Decode(&newTck); err != nil {
		log.Println("Error decoding ticket JSON:", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := tck.tsrv.CreateTicket(r.Context(), newTck)
	if err != nil {
		handleError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, map[string]any{
		"message": "Ticket created successfully",
		"id":      id,
	})
}

//
// ─── GET TICKET BY ID ────────────────────────────────────────────
func (tck *TicketsHandler) GetTicketById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "invalid ticket id", http.StatusBadRequest)
		return
	}

	ticket, err := tck.tsrv.GetTicketById(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, ticket)
}

//
// ─── DELETE TICKET ───────────────────────────────────────────────
func (tck *TicketsHandler) DeleteTicketById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "invalid ticket id", http.StatusBadRequest)
		return
	}

	err = tck.tsrv.DeleteTicketById(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Ticket deleted successfully",
	})
}

//
// ─── GET TICKETS WITH FILTERS ────────────────────────────────────
func (tck *TicketsHandler) GetTicketByFilters(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	priority := r.URL.Query().Get("priority")

	tickets, err := tck.tsrv.GetTicketsByFilters(r.Context(), status, priority)
	if err != nil {
		handleError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, tickets)
}


func (tck *TicketsHandler) UpdateTicket(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "invalid ticket id", http.StatusBadRequest)
		return
	}

	var newTck models.NewTicketFormat
	if err := json.NewDecoder(r.Body).Decode(&newTck); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = tck.tsrv.UpdateTicket(r.Context(), id, newTck)
	if err != nil {
		handleError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"message": "Ticket updated successfully",
	})
}
