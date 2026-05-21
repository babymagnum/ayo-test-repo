package request

type ReportMatchRequest struct {
	HomeScore int           `json:"home_score" binding:"required"`
	AwayScore int           `json:"away_score" binding:"required"`
	Goals     []GoalRequest `json:"goals" binding:"required"`
}

type GoalRequest struct {
	PlayerID  uint `json:"player_id" binding:"required"`
	Minute    int  `json:"minute" binding:"required"`
	IsOwnGoal bool `json:"is_own_goal"`
}
