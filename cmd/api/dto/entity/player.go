package entity

type Player struct {
	UpdateDeleteEntity
	TeamID       uint    `gorm:"not null" json:"team_id"`
	Name         string  `gorm:"type:varchar(255);not null" json:"name"`
	HeightCm     float64 `gorm:"type:numeric(5,2)" json:"height_cm"`
	WeightKg     float64 `gorm:"type:numeric(5,2)" json:"weight_kg"`
	Position     string  `gorm:"type:varchar(20);not null" json:"position"`
	JerseyNumber int     `gorm:"type:int;not null" json:"jersey_number"`
}

func (Player) TableName() string {
	return "players"
}
