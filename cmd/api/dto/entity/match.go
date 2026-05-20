package entity

import "time"

type Match struct {
	UpdateDeleteEntity
	HomeTeamID uint      `gorm:"not null" json:"home_team_id"`
	AwayTeamID uint      `gorm:"not null" json:"away_team_id"`
	MatchDate  time.Time `gorm:"type:date;not null" json:"match_date"`
	MatchTime  string    `gorm:"type:time;not null" json:"match_time"`
	Status     string    `gorm:"type:varchar(20);default:scheduled" json:"status"`
}

func (Match) TableName() string {
	return "matches"
}
