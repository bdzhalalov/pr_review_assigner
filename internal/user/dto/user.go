package dto

type GetUsersByIdsDTO struct {
	IDs []string
}

type GetUserByIDDTO struct {
	UserID string `json:"user_id"`
}

type UpdateUserActivityDTO struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type UserDTO struct {
	UserID   string
	Username string
	TeamName string
	IsActive bool
}
