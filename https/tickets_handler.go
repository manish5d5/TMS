package tickets_handler


import (
		"net/http"
		"encoding/json"
)

type TicketsHandler struct{
	tsrv *TicketsService
}

func NewTicketsHandler(s *tickets_service) *TicketsHandler{
	return &TicketsHandler{tsrv: s}
}

func (tck *TicketsHandler)NewTickets(r http.Request,w http.ResponseWriter){
	var newTck models.NewTicket_format
	json.

}