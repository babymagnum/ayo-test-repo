package interfaces

import (
	"context"

	"github.com/ariefzainuri96/ayo-test/cmd/api/dto/entity"
	"github.com/ariefzainuri96/ayo-test/cmd/api/dto/request"
	"github.com/ariefzainuri96/ayo-test/internal/utils"
)

type IPlayer interface {
	Create(ctx context.Context, teamID uint, req request.AddPlayerRequest) (entity.Player, error)
	GetByID(ctx context.Context, id uint) (entity.Player, error)
	GetByTeam(ctx context.Context, teamID uint, req request.PaginationRequest) (utils.PaginateResult[entity.Player], error)
	Update(ctx context.Context, id uint, req request.AddPlayerRequest) (entity.Player, error)
	Delete(ctx context.Context, id uint) error
}
