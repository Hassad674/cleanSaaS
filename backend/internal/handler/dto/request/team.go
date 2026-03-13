package request

type CreateTeamRequest struct {
	Name string `json:"name"`
}

type UpdateTeamRequest struct {
	Name string `json:"name"`
}

type InviteMemberRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type AcceptInviteRequest struct {
	Token string `json:"token"`
}

type DeclineInviteRequest struct {
	Token string `json:"token"`
}

type UpdateMemberRoleRequest struct {
	Role string `json:"role"`
}
