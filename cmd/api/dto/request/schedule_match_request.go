package request

type ScheduleMatchRequest struct {
	HomeTeamID uint   `json:"home_team_id" binding:"required"`
	AwayTeamID uint   `json:"away_team_id" binding:"required"`
	MatchDate  string `json:"match_date" binding:"required"`  // format: "2026-01-02"
	MatchTime  string `json:"match_time" binding:"required"`  // format: "15:04:05"
}
