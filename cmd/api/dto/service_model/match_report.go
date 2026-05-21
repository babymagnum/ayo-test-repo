package service_model

import "github.com/ariefzainuri96/go-logstream/cmd/api/dto/entity"

type MatchReport struct {
	Match              entity.Match
	HomeTeam           entity.Team
	AwayTeam           entity.Team
	Result             entity.MatchResult
	Goals              []entity.MatchGoal
	PlayerGoals        []entity.Player
	TopScorer          string
	TopScorerGoalCount int
}
