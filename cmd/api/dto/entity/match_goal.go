package entity

type MatchGoal struct {
	ID            uint `gorm:"primarykey" json:"id"`
	MatchResultID uint `gorm:"not null" json:"match_result_id"`
	PlayerID      uint `gorm:"not null" json:"player_id"`
	Minute        int  `gorm:"type:int;not null" json:"minute"`
	IsOwnGoal     bool `gorm:"type:boolean;default:false" json:"is_own_goal"`
}

func (MatchGoal) TableName() string {
	return "match_goals"
}
