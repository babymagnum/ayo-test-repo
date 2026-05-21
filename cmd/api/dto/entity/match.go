package entity

import (
	"time"
)

type Match struct {
	UpdateDeleteEntity
	HomeTeamID uint      `gorm:"not null" json:"home_team_id"`
	AwayTeamID uint      `gorm:"not null" json:"away_team_id"`
	HomeTeam   Team      `gorm:"foreignKey:HomeTeamID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"home_team"`
	AwayTeam   Team      `gorm:"foreignKey:AwayTeamID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"away_team"`
	MatchDate  time.Time `gorm:"type:date;not null" json:"match_date"`
	MatchTime  string    `gorm:"type:time;not null" json:"match_time"`
	Status     string    `gorm:"type:varchar(20);default:scheduled" json:"status"`
}

func (Match) TableName() string {
	return "matches"
}
