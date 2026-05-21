package interfaces

import (
	"context"

	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/entity"
	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/request"
	"github.com/ariefzainuri96/go-logstream/internal/utils"
)

type ITeam interface {
	Create(ctx context.Context, req request.AddTeamRequest) (entity.Team, error)
	GetByID(ctx context.Context, id uint) (entity.Team, error)
	GetAll(ctx context.Context, req request.PaginationRequest) (utils.PaginateResult[entity.Team], error)
	Update(ctx context.Context, id uint, req request.AddTeamRequest) (entity.Team, error)
	Delete(ctx context.Context, id uint) error
}
