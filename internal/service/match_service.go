package service

import (
	"context"

	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/entity"
	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/request"
	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/service_model"
	"github.com/ariefzainuri96/go-logstream/internal/store"
	"github.com/ariefzainuri96/go-logstream/internal/utils"
	"go.uber.org/zap"
)

type MatchServiceImpl struct {
	logger *zap.Logger
	store  store.Storage
}

func NewMatchService(store store.Storage, logger *zap.Logger) *MatchServiceImpl {
	return &MatchServiceImpl{
		logger: logger,
		store:  store,
	}
}

func (s *MatchServiceImpl) Schedule(ctx context.Context, req request.ScheduleMatchRequest) (entity.Match, error) {
	return s.store.IMatch.Schedule(ctx, req)
}

func (s *MatchServiceImpl) GetByID(ctx context.Context, id uint) (entity.Match, error) {
	return s.store.IMatch.GetByID(ctx, id)
}

func (s *MatchServiceImpl) GetAll(ctx context.Context, req request.PaginationRequest) (utils.PaginateResult[entity.Match], error) {
	return s.store.IMatch.GetAll(ctx, req)
}

func (s *MatchServiceImpl) Update(ctx context.Context, id uint, req request.ScheduleMatchRequest) (entity.Match, error) {
	return s.store.IMatch.Update(ctx, id, req)
}

func (s *MatchServiceImpl) Delete(ctx context.Context, id uint) error {
	return s.store.IMatch.Delete(ctx, id)
}

func (s *MatchServiceImpl) ReportResult(ctx context.Context, matchID uint, req request.ReportMatchRequest) error {
	return s.store.IMatch.ReportResult(ctx, matchID, req)
}

func (s *MatchServiceImpl) GetReport(ctx context.Context, matchID uint) (service_model.MatchReport, error) {
	match, matchResult, goals, err := s.store.IMatch.GetReport(ctx, matchID)

	if err != nil {
		return service_model.MatchReport{}, err
	}

	homeTeam, err := s.store.ITeam.GetByID(ctx, match.HomeTeamID)
	if err != nil {
		return service_model.MatchReport{}, err
	}

	awayTeam, err := s.store.ITeam.GetByID(ctx, match.AwayTeamID)
	if err != nil {
		return service_model.MatchReport{}, err
	}

	playerId, goalCount := topScorer(goals)

	player, err := s.store.IPlayer.GetByID(ctx, playerId)

	if err != nil {
		return service_model.MatchReport{}, err
	}

	distinctPlayers, err := s.store.IMatch.GetDistinctPlayers(ctx, matchResult.ID)

	if err != nil {
		return service_model.MatchReport{}, err
	}

	return service_model.MatchReport{
		Match:              match,
		Result:             matchResult,
		Goals:              goals,
		PlayerGoals:        distinctPlayers,
		HomeTeam:           homeTeam,
		AwayTeam:           awayTeam,
		TopScorerGoalCount: goalCount,
		TopScorer:          player.Name,
	}, nil
}

func topScorer(goals []entity.MatchGoal) (playerId uint, count int) {
	goalCount := make(map[uint]int)
	for _, g := range goals {
		goalCount[g.PlayerID]++
	}

	var topScorerID uint
	var maxGoal int
	for id, c := range goalCount {
		if c > maxGoal {
			topScorerID = id
			maxGoal = c
		}
	}

	return topScorerID, maxGoal
}
