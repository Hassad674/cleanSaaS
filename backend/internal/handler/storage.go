package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	appstorage "github.com/hassad/boilerplateSaaS/backend/internal/app/storage"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/dto/response"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/middleware"
)

type StorageHandler struct {
	svc *appstorage.Service
}

func NewStorageHandler(svc *appstorage.Service) *StorageHandler {
	return &StorageHandler{svc: svc}
}

func (h *StorageHandler) Upload(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	// 50MB max
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, "file too large or invalid form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "missing file field")
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	result, err := h.svc.Upload(r.Context(), userID, header.Filename, contentType, header.Size, file)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, response.FileFromDomain(result))
}

func (h *StorageHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	files, total, err := h.svc.List(r.Context(), userID, offset, limit)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	items := make([]response.FileResponse, len(files))
	for i, f := range files {
		items[i] = response.FileFromDomain(f)
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"files": items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *StorageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	fileID := chi.URLParam(r, "id")

	if err := h.svc.Delete(r.Context(), userID, fileID); err != nil {
		response.HandleDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
