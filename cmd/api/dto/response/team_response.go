package response

import "github.com/ariefzainuri96/go-logstream/cmd/api/dto/entity"

type TeamsResponse struct {
	BaseResponse
	Teams      []entity.Team    `json:"teams"`
	Pagination PaginationMetadata `json:"pagination"`
}

type TeamResponse struct {
	BaseResponse
	Team entity.Team `json:"team"`
}
