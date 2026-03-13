package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	appblog "github.com/hassad/boilerplateSaaS/backend/internal/app/blog"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/dto/request"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/dto/response"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/middleware"
)

type BlogHandler struct {
	svc *appblog.Service
}

func NewBlogHandler(svc *appblog.Service) *BlogHandler {
	return &BlogHandler{svc: svc}
}

// Public endpoints

func (h *BlogHandler) ListPublished(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Query().Get("tag")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	posts, total, err := h.svc.ListPublished(r.Context(), tag, offset, limit)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	items := make([]response.BlogPostResponse, len(posts))
	for i, p := range posts {
		items[i] = response.BlogPostFromDomain(p)
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"posts": items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *BlogHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	post, err := h.svc.GetBySlug(r.Context(), slug)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.BlogPostFromDomain(post))
}

func (h *BlogHandler) ListTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.svc.ListTags(r.Context())
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{"tags": tags})
}

// Admin endpoints

func (h *BlogHandler) AdminList(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	tag := r.URL.Query().Get("tag")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	posts, total, err := h.svc.ListAll(r.Context(), status, tag, offset, limit)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	items := make([]response.BlogPostResponse, len(posts))
	for i, p := range posts {
		items[i] = response.BlogPostFromDomain(p)
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"posts": items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *BlogHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	var req request.CreateBlogPostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	post, err := h.svc.Create(r.Context(), userID, req.Title, req.Slug, req.Excerpt, req.Content,
		req.CoverImageURL, req.MetaTitle, req.MetaDescription, req.Tags, req.Publish)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, response.BlogPostFromDomain(post))
}

func (h *BlogHandler) Update(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")

	var req request.UpdateBlogPostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	post, err := h.svc.Update(r.Context(), postID, req.Title, req.Slug, req.Excerpt, req.Content,
		req.CoverImageURL, req.MetaTitle, req.MetaDescription, req.Tags, req.Status)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.BlogPostFromDomain(post))
}

func (h *BlogHandler) Delete(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")

	if err := h.svc.Delete(r.Context(), postID); err != nil {
		response.HandleDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
