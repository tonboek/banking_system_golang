package handlers

import (
	"banking_system_golang/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type AccountHandler struct {
	service *services.AccountService
}

func NewAccountHandler(service *services.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Owner string `json:"owner"`
	}
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "ошибка при создании аккаунта", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Owner) == "" {
		http.Error(w, "владелец не должен быть пустым", http.StatusBadRequest)
		return
	}

	acc, err := h.service.CreateAccount(r.Context(), req.Owner)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(acc)
}

func (h *AccountHandler) GetAccountByID(w http.ResponseWriter, r *http.Request) {
	partsOfPath := strings.Split(r.URL.Path, "/")
	if len(partsOfPath) < 3 {
		http.Error(w, "неверный путь", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(partsOfPath[2])
	if err != nil {
		http.Error(w, "ошибка при конвертации ID", http.StatusBadRequest)
		return
	}

	acc, err := h.service.GetAccountByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(acc)
}

func (h *AccountHandler) AddMoneyToBalance(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		ID int `json:"id"`
		Amount float64 `json:"amount"`
	}
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "ошибка при пополнении баланса", http.StatusBadRequest)
		return
	}

	acc, err := h.service.AddMoneyToBalance(r.Context(), req.ID, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(acc)
}

func (h *AccountHandler) Transaction(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		FromID int `json:"from_id"`
		ToID int `json:"to_id"`
		Amount float64 `json:"amount"`
	}
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "ошибка при получении данных перевода", http.StatusBadRequest)
		return
	}

	result, err := h.service.Transaction(r.Context(), req.FromID, req.ToID, req.Amount)
	if err != nil {
		http.Error(w, "ошибка при выполнении перевода", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(result)
}
