package httphandler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mSulimenko/dev-blog-platform/internal/auth/dto"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Errorw("unable to decode", "error", err, "request", req)
		h.sendError(w, http.StatusBadRequest, ErrCodeInvalidJSON, "Invalid JSON format")
		return
	}
	defer r.Body.Close()

	if err := h.validate.Struct(req); err != nil {
		h.log.Errorw("Validation failed", "error", err)
		h.sendError(w, http.StatusBadRequest, ErrCodeValidation, "Invalid input data")
		return
	}

	userID, err := h.usersService.Register(r.Context(), &req)
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

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Errorw("unable to decode", "error", err, "request", req)
		h.sendError(w, http.StatusBadRequest, ErrCodeInvalidJSON, "Invalid JSON format")
	}
	defer r.Body.Close()

	if err := h.validate.Struct(req); err != nil {
		h.log.Errorw("Validation failed", "error", err)
		h.sendError(w, http.StatusBadRequest, ErrCodeValidation, "Invalid input data")
		return
	}

	resp, err := h.usersService.Login(r.Context(), &req)
	if err != nil {
		h.log.Errorw("Login failed", "error", err)
		h.sendError(w, http.StatusInternalServerError, ErrCodeInternal, "Server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
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
		h.sendError(w, http.StatusInternalServerError, ErrCodeInternal, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		h.sendVerificationHTML(w, false, "Verification token is required")
		return
	}

	err := h.usersService.VerifyEmail(r.Context(), token)
	if err != nil {
		h.log.Errorw("failed to verify email", "token", token, "error", err)
		h.sendVerificationHTML(w, false, "Invalid or expired verification token")
		return
	}

	h.sendVerificationHTML(w, true, "Email successfully verified!")
}

func (h *Handler) sendVerificationHTML(w http.ResponseWriter, success bool, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var html string
	if success {
		html = `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Email Verified</title>
			<style>
				body { font-family: Arial, sans-serif; text-align: center; padding: 50px; }
				.success { color: #22c55e; }
			</style>
		</head>
		<body>
			<h1 class="success">✅ Email Verified!</h1>
		</body>
		</html>`
	} else {
		html = `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Verification Failed</title>
			<style>
				body { font-family: Arial, sans-serif; text-align: center; padding: 50px; }
				.error { color: #ef4444; }
			</style>
		</head>
		<body>
			<h1 class="error">❌ Verification Failed</h1>
			<p>` + message + `</p>
			<p>Please try registering again or contact support.</p>
		</body>
		</html>`
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
