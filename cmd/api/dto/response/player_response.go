package response

import "github.com/ariefzainuri96/ayo-test/cmd/api/dto/entity"

type PlayersResponse struct {
	BaseResponse
	Players    []entity.Player  `json:"players"`
	Pagination PaginationMetadata `json:"pagination"`
}

type PlayerResponse struct {
	BaseResponse
	Player entity.Player `json:"player"`
}
