package service

import (
	"context"

	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/entity"
	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/request"
	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/service_model"
	"github.com/ariefzainuri96/go-logstream/internal/interfaces"
	"github.com/ariefzainuri96/go-logstream/internal/store"
	"github.com/ariefzainuri96/go-logstream/internal/utils"
	"go.uber.org/zap"
)

type IMatch interface {
	Schedule(ctx context.Context, req request.ScheduleMatchRequest) (entity.Match, error)
	GetByID(ctx context.Context, id uint) (entity.Match, error)
	GetAll(ctx context.Context, req request.PaginationRequest) (utils.PaginateResult[entity.Match], error)
	Update(ctx context.Context, id uint, req request.ScheduleMatchRequest) (entity.Match, error)
	Delete(ctx context.Context, id uint) error
	ReportResult(ctx context.Context, matchID uint, req request.ReportMatchRequest) error
	GetReport(ctx context.Context, matchID uint) (service_model.MatchReport, error)
}

type Service struct {
	IAuth    interfaces.IAuth
	ITeam    interfaces.ITeam
	IPlayer  interfaces.IPlayer
	IMatch   IMatch
}

func NewService(store store.Storage, logger *zap.Logger) Service {
	return Service{
		IAuth:    NewAuthService(store, logger),
		ITeam:    NewTeamService(store, logger),
		IPlayer:  NewPlayerService(store, logger),
		IMatch:   NewMatchService(store, logger),
	}
}
