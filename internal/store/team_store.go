package store

import (
	"context"
	"time"

	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/entity"
	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/request"
	db "github.com/ariefzainuri96/go-logstream/internal/db"
	"github.com/ariefzainuri96/go-logstream/internal/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TeamStore struct {
	gormDb *db.GormDB
	logger *zap.Logger
}

func (s *TeamStore) Create(ctx context.Context, req request.AddTeamRequest) (entity.Team, error) {
	team := entity.Team{
		Name:        req.Name,
		LogoURL:     req.LogoURL,
		FoundedYear: req.FoundedYear,
		Address:     req.Address,
		City:        req.City,
	}

	err := s.gormDb.ExecWithTimeoutErr(ctx, func(tx *gorm.DB) error {
		return tx.Create(&team).Error
	})

	if err != nil {
		return entity.Team{}, err
	}

	return team, nil
}

func (s *TeamStore) GetByID(ctx context.Context, id uint) (entity.Team, error) {
	var team entity.Team

	err := s.gormDb.ExecWithTimeoutErr(ctx, func(tx *gorm.DB) error {
		return tx.First(&team, id).Error
	})

	if err != nil {
		return entity.Team{}, err
	}

	return team, nil
}

func (s *TeamStore) GetAll(ctx context.Context, req request.PaginationRequest) (utils.PaginateResult[entity.Team], error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	query := s.gormDb.GormDb.WithContext(ctx).Model(&entity.Team{})

	result := utils.ApplyPagination[entity.Team](query, req, "")

	if result.Error != nil {
		return utils.PaginateResult[entity.Team]{}, result.Error
	}

	return result, nil
}

func (s *TeamStore) Update(ctx context.Context, id uint, req request.AddTeamRequest) (entity.Team, error) {
	var team entity.Team

	err := s.gormDb.ExecWithTimeoutErr(ctx, func(tx *gorm.DB) error {
		return tx.First(&team, id).Error
	})

	if err != nil {
		return entity.Team{}, err
	}

	team.Name = req.Name
	team.LogoURL = req.LogoURL
	team.FoundedYear = req.FoundedYear
	team.Address = req.Address
	team.City = req.City

	err = s.gormDb.ExecWithTimeoutErr(ctx, func(tx *gorm.DB) error {
		return tx.Save(&team).Error
	})

	if err != nil {
		return entity.Team{}, err
	}

	return team, nil
}

func (s *TeamStore) Delete(ctx context.Context, id uint) error {
	return s.gormDb.ExecWithTimeoutErr(ctx, func(tx *gorm.DB) error {
		return tx.Delete(&entity.Team{}, id).Error
	})
}
