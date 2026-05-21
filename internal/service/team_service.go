package service

import (
	"context"

	"github.com/ariefzainuri96/ayo-test/cmd/api/dto/entity"
	"github.com/ariefzainuri96/ayo-test/cmd/api/dto/request"
	"github.com/ariefzainuri96/ayo-test/internal/store"
	"github.com/ariefzainuri96/ayo-test/internal/utils"
	"go.uber.org/zap"
)

type TeamServiceImpl struct {
	logger *zap.Logger
	store  store.Storage
}

func NewTeamService(store store.Storage, logger *zap.Logger) *TeamServiceImpl {
	return &TeamServiceImpl{
		logger: logger,
		store:  store,
	}
}

func (s *TeamServiceImpl) Create(ctx context.Context, req request.AddTeamRequest) (entity.Team, error) {
	return s.store.ITeam.Create(ctx, req)
}

func (s *TeamServiceImpl) GetByID(ctx context.Context, id uint) (entity.Team, error) {
	return s.store.ITeam.GetByID(ctx, id)
}

func (s *TeamServiceImpl) GetAll(ctx context.Context, req request.PaginationRequest) (utils.PaginateResult[entity.Team], error) {
	return s.store.ITeam.GetAll(ctx, req)
}

func (s *TeamServiceImpl) Update(ctx context.Context, id uint, req request.AddTeamRequest) (entity.Team, error) {
	return s.store.ITeam.Update(ctx, id, req)
}

func (s *TeamServiceImpl) Delete(ctx context.Context, id uint) error {
	return s.store.ITeam.Delete(ctx, id)
}
