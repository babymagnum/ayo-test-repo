package store

import (
	"context"

	"github.com/ariefzainuri96/ayo-test/cmd/api/dto/entity"
	"github.com/ariefzainuri96/ayo-test/cmd/api/dto/request"
	db "github.com/ariefzainuri96/ayo-test/internal/db"
	"github.com/ariefzainuri96/ayo-test/internal/interfaces"
	"github.com/ariefzainuri96/ayo-test/internal/utils"
	"go.uber.org/zap"
)

type IMatch interface {
	Schedule(ctx context.Context, req request.ScheduleMatchRequest) (entity.Match, error)
	GetByID(ctx context.Context, id uint) (entity.Match, error)
	GetAll(ctx context.Context, req request.PaginationRequest) (utils.PaginateResult[entity.Match], error)
	Update(ctx context.Context, id uint, req request.ScheduleMatchRequest) (entity.Match, error)
	Delete(ctx context.Context, id uint) error
	ReportResult(ctx context.Context, matchID uint, req request.ReportMatchRequest) error
	GetReport(ctx context.Context, matchID uint) (entity.Match, entity.MatchResult, []entity.MatchGoal, error)
	GetDistinctPlayers(ctx context.Context, matchResultId uint) ([]entity.Player, error)
}

type Storage struct {
	IAuth   interfaces.IAuth
	ITeam   interfaces.ITeam
	IPlayer interfaces.IPlayer
	IMatch  IMatch
}

func NewStorage(gorm *db.GormDB, logger *zap.Logger) Storage {
	return Storage{
		IAuth:   &AuthStore{gorm},
		ITeam:   &TeamStore{gorm, logger},
		IPlayer: &PlayerStore{gorm, logger},
		IMatch:  &MatchStore{gorm, logger},
	}
}
