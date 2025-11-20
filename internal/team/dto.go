package team

type AddTeamDTO struct {
	TeamName string
	Members  []struct {
		UserID   string
		Username string
		IsActive bool
	}
}
