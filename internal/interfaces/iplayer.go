package interfaces

import (
	"context"

	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/entity"
	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/request"
	"github.com/ariefzainuri96/go-logstream/internal/utils"
)

type IPlayer interface {
	Create(ctx context.Context, teamID uint, req request.AddPlayerRequest) (entity.Player, error)
	GetByID(ctx context.Context, id uint) (entity.Player, error)
	GetByTeam(ctx context.Context, teamID uint, req request.PaginationRequest) (utils.PaginateResult[entity.Player], error)
	Update(ctx context.Context, id uint, req request.AddPlayerRequest) (entity.Player, error)
	Delete(ctx context.Context, id uint) error
}
