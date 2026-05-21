package response

import (
	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/entity"
)

type MatchGoalData struct {
	ID         uint   `json:"id"`
	PlayerID   uint   `json:"player_id"`
	PlayerName string `json:"player_name"`
	Minute     int    `json:"minute"`
	IsOwnGoal  bool   `json:"is_own_goal"`
}

type MatchReportResponse struct {
	BaseResponse
	Data MatchReportData `json:"data"`
}

type MatchReportData struct {
	Match              entity.Match    `json:"match"`
	HomeTeam           entity.Team     `json:"home_team"`
	AwayTeam           entity.Team     `json:"away_team"`
	HomeScore          int             `json:"home_score"`
	AwayScore          int             `json:"away_score"`
	Winner             string          `json:"winner"`
	Goals              []MatchGoalData `json:"goals"`
	TopScorer          string          `json:"top_scorer"`
	TopScorerGoalCount int             `json:"top_scorer_goal_count"`
	HomeTotalWins      int             `json:"home_total_wins"`
	AwayTotalWins      int             `json:"away_total_wins"`
}
