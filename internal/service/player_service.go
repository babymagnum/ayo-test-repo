package service

import (
	"context"
	"errors"
	"strings"

	"github.com/ariefzainuri96/ayo-test/cmd/api/dto/entity"
	"github.com/ariefzainuri96/ayo-test/cmd/api/dto/request"
	"github.com/ariefzainuri96/ayo-test/internal/store"
	"github.com/ariefzainuri96/ayo-test/internal/utils"
	"go.uber.org/zap"
)

var validPositions = map[string]bool{
	"penyerang":        true,
	"gelandang":        true,
	"bertahan":         true,
	"penjaga gawang":   true,
}

func validatePosition(position string) error {
	if !validPositions[strings.ToLower(position)] {
		return errors.New("posisi harus salah satu dari: penyerang, gelandang, bertahan, penjaga gawang")
	}
	return nil
}

type PlayerServiceImpl struct {
	logger *zap.Logger
	store  store.Storage
}

func NewPlayerService(store store.Storage, logger *zap.Logger) *PlayerServiceImpl {
	return &PlayerServiceImpl{
		logger: logger,
		store:  store,
	}
}

func (s *PlayerServiceImpl) Create(ctx context.Context, teamID uint, req request.AddPlayerRequest) (entity.Player, error) {
	if err := validatePosition(req.Position); err != nil {
		return entity.Player{}, err
	}
	return s.store.IPlayer.Create(ctx, teamID, req)
}

func (s *PlayerServiceImpl) GetByID(ctx context.Context, id uint) (entity.Player, error) {
	return s.store.IPlayer.GetByID(ctx, id)
}

func (s *PlayerServiceImpl) GetByTeam(ctx context.Context, teamID uint, req request.PaginationRequest) (utils.PaginateResult[entity.Player], error) {
	return s.store.IPlayer.GetByTeam(ctx, teamID, req)
}

func (s *PlayerServiceImpl) Update(ctx context.Context, id uint, req request.AddPlayerRequest) (entity.Player, error) {
	if err := validatePosition(req.Position); err != nil {
		return entity.Player{}, err
	}
	return s.store.IPlayer.Update(ctx, id, req)
}

func (s *PlayerServiceImpl) Delete(ctx context.Context, id uint) error {
	return s.store.IPlayer.Delete(ctx, id)
}
