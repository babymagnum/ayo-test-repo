package entity

type Team struct {
	UpdateDeleteEntity
	Name        string `gorm:"unique;type:varchar(255);not null" json:"name"`
	LogoURL     string `gorm:"type:text" json:"logo_url"`
	FoundedYear int    `gorm:"type:int" json:"founded_year"`
	Address     string `gorm:"type:text" json:"address"`
	City        string `gorm:"type:varchar(100)" json:"city"`
	TotalWins   int    `gorm:"type:int;default:0" json:"total_wins"`
}

func (Team) TableName() string {
	return "teams"
}
