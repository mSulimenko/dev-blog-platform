package httphandler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mSulimenko/dev-blog-platform/internal/articles/dto"
)

func (h *Handler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Errorw("Failed to decode request body", "error", err)
		sendError(w, http.StatusBadRequest, ErrCodeInvalidJSON, "Invalid JSON format")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		h.log.Warnw("Validation failed", "error", err)
		sendError(w, http.StatusBadRequest, ErrCodeValidation, "Invalid request")
		return
	}

	userID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		sendError(w, http.StatusUnauthorized, ErrCodeValidation, "Authentication required")
		return
	}

	if req.AuthorId != userID {
		sendError(w, http.StatusForbidden, ErrCodeValidation, "Can only create articles for yourself")
		return
	}

	article, err := h.articlesService.CreateArticle(r.Context(), req)
	if err != nil {
		h.log.Errorw("Failed to create article", "error", err)
		http.Error(w, `{"error": "Failed to create article"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(article)
}

func (h *Handler) GetArticle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		sendError(w, http.StatusBadRequest, ErrCodeValidation, "id required")
		return
	}

	article, err := h.articlesService.GetArticle(r.Context(), id)
	if err != nil {
		h.log.Errorw("Failed to get article", "id", id, "error", err)
		sendError(w, http.StatusBadRequest, ErrCodeNotFound, "Article not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)
}

func (h *Handler) ListArticles(w http.ResponseWriter, r *http.Request) {
	var req dto.ListRequest
	query := r.URL.Query()

	if authorID := query.Get("author_id"); authorID != "" {
		req.AuthorID = &authorID
	}
	if status := query.Get("status"); status != "" {
		req.Status = &status
	}
	offsetStr := query.Get("offset")
	if offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			req.Offset = offset
		}
	}
	limitStr := query.Get("limit")
	if limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = limit
		}
	}
	if req.Limit == 0 {
		req.Limit = 20
	}

	if err := h.validate.Struct(req); err != nil {
		h.log.Warnw("Validation failed for list request", "error", err)
		sendError(w, http.StatusBadRequest, ErrCodeValidation, "Invalid parameters")
		return
	}

	articles, err := h.articlesService.ListArticles(r.Context(), req)
	if err != nil {
		h.log.Errorw("Failed to list articles", "error", err)
		sendError(w, http.StatusInternalServerError, ErrCodeInternal, "internal error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articles)
}

func (h *Handler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	articleId := chi.URLParam(r, "id")
	if articleId == "" {
		sendError(w, http.StatusBadRequest, ErrCodeValidation, "id required")
		return
	}

	var req dto.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Errorw("Failed to decode request body", "error", err)
		sendError(w, http.StatusBadRequest, ErrCodeInvalidJSON, "Invalid JSON format")
		return
	}

	if req.Title == nil && req.Content == nil && req.Status == nil {
		sendError(w, http.StatusBadRequest, ErrCodeValidation, "At least one field must be provided for update")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		h.log.Warnw("Validation failed for update request", "error", err)
		sendError(w, http.StatusBadRequest, ErrCodeValidation, "Invalid request")
		return
	}
	userID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		sendError(w, http.StatusUnauthorized, ErrCodeValidation, "Authentication required")
		return
	}

	article, err := h.articlesService.UpdateArticle(r.Context(), articleId, req, userID)
	if err != nil {
		h.log.Errorw("Failed to update article", "id", articleId, "error", err)
		sendError(w, http.StatusInternalServerError, ErrCodeInternal, "Internal error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)
}

func (h *Handler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	articleId := chi.URLParam(r, "id")
	if articleId == "" {
		sendError(w, http.StatusBadRequest, ErrCodeValidation, "id required")
		return
	}
	userID, err := h.getUserIDFromContext(r.Context())
	if err != nil {
		sendError(w, http.StatusUnauthorized, ErrCodeValidation, "Authentication required")
		return
	}
	err = h.articlesService.DeleteArticle(r.Context(), articleId, userID)
	if err != nil {
		h.log.Errorw("Failed to delete article", "id", articleId, "error", err)
		sendError(w, http.StatusInternalServerError, ErrCodeInternal, "Internal error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

func (h *Handler) getUserRoleFromContext(ctx context.Context) (string, error) {
	role, ok := ctx.Value("user_role").(string)
	if !ok || role == "" {
		return "", fmt.Errorf("user role not found in context")
	}
	return role, nil
}
