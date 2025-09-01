package handlers

import (
	"banking_system_golang/internal/services"
	"encoding/json"
	"net/http"
	"strings"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func(h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.Register(r.Context(), req.Name, req.Username, req.Email, req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			http.Error(w, "username или email уже заняты", http.StatusConflict)
			return
		}
		http.Error(w, "ошибка при регистрации", http.StatusBadRequest)
		return
	}

	user.Password = ""
	json.NewEncoder(w).Encode(&models.AuthResponse{Token: token, User: *user})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	user, err := h.service.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(&models.AuthResponse{Token: token, User: *user})
}