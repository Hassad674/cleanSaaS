package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	appstorage "github.com/hassad/boilerplateSaaS/backend/internal/app/storage"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/dto/response"
	"github.com/hassad/boilerplateSaaS/backend/pkg/orgctx"
)

// demoUserID is a fixed UUID used for all demo storage operations.
// This keeps demo files isolated from real user data.
const demoUserID = "00000000-0000-0000-0000-000000000000"

// demoOrgID is the fixed organization the public demo runs under. The demo is a
// real tenant: its data is org-scoped and RLS-isolated like any other org, so the
// public demos keep working end-to-end without seeing (or leaking into) real
// tenants. It is seeded by cmd/seed alongside the demo user.
const demoOrgID = "00000000-0000-0000-0000-0000000000d0"

// withDemoOrg scopes a request context to the fixed demo organization so the
// org-scoped storage path (and its row-level security) applies to demo files.
func withDemoOrg(ctx context.Context) context.Context {
	return orgctx.WithOrgID(ctx, demoOrgID)
}

// DemoStorageHandler provides public (no-auth) storage endpoints for the
// landing-page demo. It reuses the real storage app service but scopes
// every operation to a hardcoded demo user ID.
type DemoStorageHandler struct {
	svc *appstorage.Service
}

func NewDemoStorageHandler(svc *appstorage.Service) *DemoStorageHandler {
	return &DemoStorageHandler{svc: svc}
}

// Upload handles POST /demo/storage/upload
// Accepts a multipart file upload (max 10MB for demo) and stores it in R2.
func (h *DemoStorageHandler) Upload(w http.ResponseWriter, r *http.Request) {
	// 10MB max for demo uploads
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, "file too large or invalid form (max 10MB for demo)")
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

	result, err := h.svc.Upload(withDemoOrg(r.Context()), demoUserID, header.Filename, contentType, header.Size, file)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, response.FileFromDomain(result))
}

// List handles GET /demo/storage/files
// Returns paginated list of demo files.
func (h *DemoStorageHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	files, total, err := h.svc.List(withDemoOrg(r.Context()), demoUserID, offset, limit)
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

// Delete handles DELETE /demo/storage/files/{id}
// Deletes a demo file from R2 and the database.
func (h *DemoStorageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	fileID := chi.URLParam(r, "id")

	if err := h.svc.Delete(withDemoOrg(r.Context()), demoUserID, fileID); err != nil {
		response.HandleDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
