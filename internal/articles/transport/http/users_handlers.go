package httphandler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mSulimenko/dev-blog-platform/internal/articles/dto"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Errorw("unable to decode", "error", err, "request", req)
		h.sendError(w, http.StatusBadRequest, ErrCodeInvalidJSON, "Invalid JSON format")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		h.log.Errorw("Validation failed", "error", err)
		h.sendError(w, http.StatusBadRequest, ErrCodeValidation, "Invalid input data")
		return
	}

	userID, err := h.usersService.CreateUser(r.Context(), &req)
	if err != nil {
		h.log.Errorw("Failed to create user", "error", err, "request", req)
		h.sendError(w, http.StatusInternalServerError, ErrCodeInternal, "Internal server error")
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
		h.sendError(w, http.StatusInternalServerError, ErrCodeInternal, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Errorw("Failed to decode request", "error", err)
		h.sendError(w, http.StatusBadRequest, ErrCodeValidation, "Invalid JSON format")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		h.log.Errorw("Failed to validate request", "error", err)
		h.sendError(w, http.StatusBadRequest, ErrCodeValidation, "Invalid input data")
		return
	}

	err := h.usersService.UpdateUser(r.Context(), id, &req)
	if err != nil {
		h.log.Errorw("Failed to update user", "id", id, "error", err)
		h.sendError(w, http.StatusInternalServerError, ErrCodeInternal, "Internal server error")
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
