package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/entity"
	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/request"
	db "github.com/ariefzainuri96/go-logstream/internal/db"
	"github.com/ariefzainuri96/go-logstream/internal/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MatchStore struct {
	gormDb *db.GormDB
	logger *zap.Logger
}

func (s *MatchStore) Schedule(ctx context.Context, req request.ScheduleMatchRequest) (entity.Match, error) {
	match := entity.Match{
		HomeTeamID: req.HomeTeamID,
		AwayTeamID: req.AwayTeamID,
		MatchDate:  parseDate(req.MatchDate),
		MatchTime:  req.MatchTime,
		Status:     "scheduled",
	}

	err := s.gormDb.ExecWithTimeoutErr(ctx, func(tx *gorm.DB) error {
		return tx.Create(&match).Error
	})

	if err != nil {
		return entity.Match{}, err
	}

	return match, nil
}

func (s *MatchStore) GetByID(ctx context.Context, id uint) (entity.Match, error) {
	var match entity.Match

	err := s.gormDb.ExecWithTimeoutErr(ctx, func(tx *gorm.DB) error {
		return tx.Preload("HomeTeam").Preload("AwayTeam").First(&match, id).Error
	})

	if err != nil {
		return entity.Match{}, err
	}

	return match, nil
}

func (s *MatchStore) GetAll(ctx context.Context, req request.PaginationRequest) (utils.PaginateResult[entity.Match], error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	query := s.gormDb.GormDb.WithContext(ctx).
		Model(&entity.Match{}).
		Preload("HomeTeam").
		Preload("AwayTeam")

	result := utils.ApplyPagination[entity.Match](query, req, "")

	if result.Error != nil {
		return utils.PaginateResult[entity.Match]{}, result.Error
	}

	return result, nil
}

func (s *MatchStore) Update(ctx context.Context, id uint, req request.ScheduleMatchRequest) (entity.Match, error) {
	var match entity.Match

	err := s.gormDb.ExecWithTimeoutErr(ctx, func(tx *gorm.DB) error {
		return tx.First(&match, id).Error
	})

	if err != nil {
		return entity.Match{}, err
	}

	match.HomeTeamID = req.HomeTeamID
	match.AwayTeamID = req.AwayTeamID
	match.MatchDate = parseDate(req.MatchDate)
	match.MatchTime = req.MatchTime

	err = s.gormDb.ExecWithTimeoutErr(ctx, func(tx *gorm.DB) error {
		return tx.Save(&match).Error
	})

	if err != nil {
		return entity.Match{}, err
	}

	return match, nil
}

func (s *MatchStore) Delete(ctx context.Context, id uint) error {
	return s.gormDb.ExecWithTimeoutErr(ctx, func(tx *gorm.DB) error {
		return tx.Delete(&entity.Match{}, id).Error
	})
}

func (s *MatchStore) ReportResult(ctx context.Context, matchID uint, req request.ReportMatchRequest) error {
	return s.gormDb.ExecWithTimeoutErr(ctx, func(tx *gorm.DB) error {
		var match entity.Match
		if err := tx.First(&match, matchID).Error; err != nil {
			return errors.New("pertandingan tidak ditemukan")
		}

		if match.Status == "completed" {
			return errors.New("pertandingan sudah memiliki hasil")
		}

		winner := "draw"
		if req.HomeScore > req.AwayScore {
			winner = "home"
		} else if req.AwayScore > req.HomeScore {
			winner = "away"
		}

		result := entity.MatchResult{
			MatchID:   matchID,
			HomeScore: req.HomeScore,
			AwayScore: req.AwayScore,
			Winner:    winner,
		}

		if err := tx.Create(&result).Error; err != nil {
			return err
		}

		for _, g := range req.Goals {
			goal := entity.MatchGoal{
				MatchResultID: result.ID,
				PlayerID:      g.PlayerID,
				Minute:        g.Minute,
				IsOwnGoal:     g.IsOwnGoal,
			}
			if err := tx.Create(&goal).Error; err != nil {
				return fmt.Errorf("gagal menyimpan gol: %w", err)
			}
		}

		if err := tx.Model(&match).Update("status", "completed").Error; err != nil {
			return err
		}

		if winner != "draw" {
			var winnerTeamID uint
			if winner == "home" {
				winnerTeamID = match.HomeTeamID
			} else {
				winnerTeamID = match.AwayTeamID
			}

			if err := tx.Model(&entity.Team{}).
				Where("id = ?", winnerTeamID).
				Update("total_wins", gorm.Expr("total_wins + 1")).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *MatchStore) GetReport(ctx context.Context, matchID uint) (entity.Match, entity.MatchResult, []entity.MatchGoal, error) {
	var match entity.Match
	var result entity.MatchResult
	var goals []entity.MatchGoal

	err := s.gormDb.ExecWithTimeoutErr(ctx, func(tx *gorm.DB) error {
		if err := tx.First(&match, matchID).Error; err != nil {
			return errors.New("pertandingan tidak ditemukan")
		}

		if err := tx.Where("match_id = ?", matchID).First(&result).Error; err != nil {
			return errors.New("hasil pertandingan belum tersedia")
		}

		if err := tx.Where("match_result_id = ?", result.ID).Find(&goals).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return entity.Match{}, entity.MatchResult{}, nil, err
	}

	return match, result, goals, nil
}

func (s *MatchStore) GetDistinctPlayers(ctx context.Context, matchResultId uint) ([]entity.Player, error) {
	var goals []entity.MatchGoal

	err := s.gormDb.ExecWithTimeoutErr(ctx, func(tx *gorm.DB) error {
		return tx.Select("DISTINCT ON (player_id) *").Where("match_result_id = ?", matchResultId).Find(&goals).Error
	})

	if err != nil {
		return nil, err
	}

	playerIDs := make([]uint, 0, len(goals))
	for _, g := range goals {
		playerIDs = append(playerIDs, g.PlayerID)
	}

	var players []entity.Player
	err = s.gormDb.ExecWithTimeoutErr(ctx, func(tx *gorm.DB) error {
		return tx.Where("id IN ?", playerIDs).Find(&players).Error
	})

	if err != nil {
		return nil, err
	}

	return players, nil
}

func parseDate(dateStr string) time.Time {
	t, _ := time.Parse("2006-01-02", dateStr)
	return t
}
