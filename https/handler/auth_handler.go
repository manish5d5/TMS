package handler

import (
	"TMS/models"
	service "TMS/services"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	AuthService *service.AuthService
	UserService *service.UserService
}

func NewAuthHandler(authService *service.AuthService, userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
		UserService: userService,
	}
}
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := user.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userCreated, err := h.UserService.CreateUser(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(userCreated)
	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var userReq models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := userReq.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, cookiesModel, err := h.AuthService.Login(r.Context(), userReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if cookiesModel != nil {
		cookiesModel.ChangeSameSiteForDevelopment(r)
		http.SetCookie(w, cookiesModel.AccessCookie)
		http.SetCookie(w, cookiesModel.RefreshCookie)
	}
	json.NewEncoder(w).Encode(user)
	w.WriteHeader(http.StatusOK)

}
