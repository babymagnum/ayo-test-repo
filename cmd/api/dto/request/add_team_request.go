package request

type AddTeamRequest struct {
	Name        string `json:"name" binding:"required,max=255"`
	LogoURL     string `json:"logo_url"`
	FoundedYear int    `json:"founded_year"`
	Address     string `json:"address"`
	City        string `json:"city" binding:"max=100"`
}
