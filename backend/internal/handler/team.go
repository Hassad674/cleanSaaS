package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	appteam "github.com/hassad/boilerplateSaaS/backend/internal/app/team"
	domainteam "github.com/hassad/boilerplateSaaS/backend/internal/domain/team"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/dto/request"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/dto/response"
	"github.com/hassad/boilerplateSaaS/backend/internal/handler/middleware"
)

type TeamHandler struct {
	svc *appteam.Service
}

func NewTeamHandler(svc *appteam.Service) *TeamHandler {
	return &TeamHandler{svc: svc}
}

func (h *TeamHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	var req request.CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		response.Error(w, http.StatusBadRequest, "team name is required")
		return
	}

	t, err := h.svc.CreateTeam(r.Context(), userID, req.Name)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, response.TeamFromDomain(t))
}

func (h *TeamHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	teams, err := h.svc.ListUserTeams(r.Context(), userID)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	items := make([]response.TeamResponse, len(teams))
	for i, t := range teams {
		items[i] = response.TeamFromDomain(t)
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"teams": items,
	})
}

func (h *TeamHandler) Get(w http.ResponseWriter, r *http.Request) {
	teamID := chi.URLParam(r, "id")

	t, err := h.svc.GetTeam(r.Context(), teamID)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.TeamFromDomain(t))
}

func (h *TeamHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	teamID := chi.URLParam(r, "id")

	var req request.UpdateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		response.Error(w, http.StatusBadRequest, "team name is required")
		return
	}

	t, err := h.svc.UpdateTeam(r.Context(), userID, teamID, req.Name)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.TeamFromDomain(t))
}

func (h *TeamHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	teamID := chi.URLParam(r, "id")

	if err := h.svc.DeleteTeam(r.Context(), userID, teamID); err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "team deleted successfully"})
}

func (h *TeamHandler) InviteMember(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	teamID := chi.URLParam(r, "id")

	var req request.InviteMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" {
		response.Error(w, http.StatusBadRequest, "email is required")
		return
	}
	if req.Role == "" {
		response.Error(w, http.StatusBadRequest, "role is required")
		return
	}
	if !domainteam.ValidateRole(req.Role) {
		response.Error(w, http.StatusBadRequest, "invalid role")
		return
	}

	member, err := h.svc.InviteMember(r.Context(), userID, teamID, req.Email, domainteam.Role(req.Role))
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, response.TeamMemberFromDomain(member))
}

func (h *TeamHandler) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())

	var req request.AcceptInviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Token == "" {
		response.Error(w, http.StatusBadRequest, "invite token is required")
		return
	}

	member, err := h.svc.AcceptInvite(r.Context(), userID, req.Token)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.TeamMemberFromDomain(member))
}

func (h *TeamHandler) DeclineInvite(w http.ResponseWriter, r *http.Request) {
	var req request.DeclineInviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Token == "" {
		response.Error(w, http.StatusBadRequest, "invite token is required")
		return
	}

	if err := h.svc.DeclineInvite(r.Context(), req.Token); err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "invite declined"})
}

func (h *TeamHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	teamID := chi.URLParam(r, "id")

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	members, total, err := h.svc.ListMembers(r.Context(), userID, teamID, offset, limit)
	if err != nil {
		response.HandleDomainError(w, err)
		return
	}

	items := make([]response.TeamMemberResponse, len(members))
	for i, m := range members {
		items[i] = response.TeamMemberFromDomain(m)
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"members": items,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

func (h *TeamHandler) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	teamID := chi.URLParam(r, "id")
	targetUserID := chi.URLParam(r, "userId")

	var req request.UpdateMemberRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Role == "" {
		response.Error(w, http.StatusBadRequest, "role is required")
		return
	}
	if !domainteam.ValidateRole(req.Role) {
		response.Error(w, http.StatusBadRequest, "invalid role")
		return
	}

	if err := h.svc.UpdateMemberRole(r.Context(), userID, teamID, targetUserID, domainteam.Role(req.Role)); err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "member role updated"})
}

func (h *TeamHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	teamID := chi.URLParam(r, "id")
	targetUserID := chi.URLParam(r, "userId")

	if err := h.svc.RemoveMember(r.Context(), userID, teamID, targetUserID); err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "member removed"})
}

func (h *TeamHandler) Leave(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromContext(r.Context())
	teamID := chi.URLParam(r, "id")

	if err := h.svc.LeaveTeam(r.Context(), userID, teamID); err != nil {
		response.HandleDomainError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "left team successfully"})
}
