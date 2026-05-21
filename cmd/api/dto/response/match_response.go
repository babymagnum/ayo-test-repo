package response

import (
	"time"

	"github.com/ariefzainuri96/ayo-test/cmd/api/dto/entity"
)

type MatchesResponse struct {
	BaseResponse
	Matches    []MatchData        `json:"matches"`
	Pagination PaginationMetadata `json:"pagination"`
}

type MatchResponse struct {
	BaseResponse
	Match entity.Match `json:"match"`
}

type MatchData struct {
	ID        uint      `json:"id"`
	MatchDate time.Time `gorm:"type:date;not null" json:"match_date"`
	MatchTime string    `gorm:"type:time;not null" json:"match_time"`
	Status    string    `gorm:"type:varchar(20);default:scheduled" json:"status"`
	HomeTeam  TeamData  `json:"home_team"`
	AwayTeam  TeamData  `json:"away_team"`
}

type TeamData struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}
