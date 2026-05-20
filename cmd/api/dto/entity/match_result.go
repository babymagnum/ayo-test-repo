package entity

type MatchResult struct {
	UpdateDeleteEntity
	MatchID   uint   `gorm:"unique;not null" json:"match_id"`
	HomeScore int    `gorm:"type:int;default:0;not null" json:"home_score"`
	AwayScore int    `gorm:"type:int;default:0;not null" json:"away_score"`
	Winner    string `gorm:"type:varchar(10);not null" json:"winner"`
}

func (MatchResult) TableName() string {
	return "match_results"
}
