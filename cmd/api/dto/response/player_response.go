package response

import "github.com/ariefzainuri96/go-logstream/cmd/api/dto/entity"

type PlayersResponse struct {
	BaseResponse
	Players    []entity.Player  `json:"players"`
	Pagination PaginationMetadata `json:"pagination"`
}

type PlayerResponse struct {
	BaseResponse
	Player entity.Player `json:"player"`
}
