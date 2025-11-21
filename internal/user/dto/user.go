package dto

type GetUsersByIdsDTO struct {
	IDs []string
}

type UserDTO struct {
	UserID   string
	Username string
	TeamName string
	IsActive bool
}
