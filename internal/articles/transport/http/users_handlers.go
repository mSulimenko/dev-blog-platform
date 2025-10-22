package httphandler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mSulimenko/dev-blog-platform/internal/articles/dto"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.UserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Errorw("Failed to decode request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := h.usersService.CreateUser(r.Context(), &req)
	if err != nil {
		h.log.Errorw("Failed to create user", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id": userID,
	})
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := h.usersService.GetUser(r.Context(), id)
	if err != nil {
		h.log.Errorw("Failed to get user", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.usersService.ListUsers(r.Context())
	if err != nil {
		h.log.Errorw("Failed to list users", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.UserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Errorw("Failed to decode request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.usersService.UpdateUser(r.Context(), id, &req)
	if err != nil {
		h.log.Errorw("Failed to update user", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.usersService.DeleteUser(r.Context(), id)
	if err != nil {
		h.log.Errorw("Failed to delete user", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
