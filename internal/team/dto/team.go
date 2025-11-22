package dto

type TeamRequestDTO struct {
	TeamName string
}

type MemberDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type TeamResponseDTO struct {
	ID       uint        `json:"id"`
	TeamName string      `json:"team_name"`
	Members  []MemberDTO `json:"members"`
}
