package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/ariefzainuri96/ayo-test/cmd/api/dto/entity"
	"github.com/ariefzainuri96/ayo-test/cmd/api/dto/request"
	db "github.com/ariefzainuri96/ayo-test/internal/db"
	"github.com/ariefzainuri96/ayo-test/internal/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PlayerStore struct {
	gormDb *db.GormDB
	logger *zap.Logger
}

func (s *PlayerStore) Create(ctx context.Context, teamID uint, req request.AddPlayerRequest) (entity.Player, error) {
	player := entity.Player{
		TeamID:       teamID,
		Name:         req.Name,
		HeightCm:     req.HeightCm,
		WeightKg:     req.WeightKg,
		Position:     req.Position,
		JerseyNumber: req.JerseyNumber,
	}

	err := s.gormDb.ExecWithTimeoutErr(ctx, func(tx *gorm.DB) error {
		return tx.Create(&player).Error
	})

	if err != nil {
		if isDuplicateJerseyError(err) {
			return entity.Player{}, errors.New("nomor punggung sudah digunakan dalam tim ini")
		}
		return entity.Player{}, err
	}

	return player, nil
}

func (s *PlayerStore) GetByID(ctx context.Context, id uint) (entity.Player, error) {
	var player entity.Player

	err := s.gormDb.ExecWithTimeoutErr(ctx, func(tx *gorm.DB) error {
		return tx.First(&player, id).Error
	})

	if err != nil {
		return entity.Player{}, err
	}

	return player, nil
}

func (s *PlayerStore) GetByTeam(ctx context.Context, teamID uint, req request.PaginationRequest) (utils.PaginateResult[entity.Player], error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	query := s.gormDb.GormDb.WithContext(ctx).
		Model(&entity.Player{}).
		Where("team_id = ?", teamID)

	result := utils.ApplyPagination[entity.Player](query, req, "")

	if result.Error != nil {
		return utils.PaginateResult[entity.Player]{}, result.Error
	}

	return result, nil
}

func (s *PlayerStore) Update(ctx context.Context, id uint, req request.AddPlayerRequest) (entity.Player, error) {
	var player entity.Player

	err := s.gormDb.ExecWithTimeoutErr(ctx, func(tx *gorm.DB) error {
		return tx.First(&player, id).Error
	})

	if err != nil {
		return entity.Player{}, err
	}

	player.Name = req.Name
	player.HeightCm = req.HeightCm
	player.WeightKg = req.WeightKg
	player.Position = req.Position
	player.JerseyNumber = req.JerseyNumber

	err = s.gormDb.ExecWithTimeoutErr(ctx, func(tx *gorm.DB) error {
		return tx.Save(&player).Error
	})

	if err != nil {
		if isDuplicateJerseyError(err) {
			return entity.Player{}, errors.New("nomor punggung sudah digunakan dalam tim ini")
		}
		return entity.Player{}, err
	}

	return player, nil
}

func (s *PlayerStore) Delete(ctx context.Context, id uint) error {
	return s.gormDb.ExecWithTimeoutErr(ctx, func(tx *gorm.DB) error {
		return tx.Delete(&entity.Player{}, id).Error
	})
}

func isDuplicateJerseyError(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "players_team_id_jersey_number_key"))
}
