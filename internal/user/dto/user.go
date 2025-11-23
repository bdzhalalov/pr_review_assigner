package dto

type GetUserByIDDTO struct {
	UserID string
}

type UpdateUserActivityDTO struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type UserRequestDTO struct {
	UserID   string
	Username string
	TeamID   uint
	IsActive bool
}

type UserResponseDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}
