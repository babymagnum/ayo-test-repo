package response

import "github.com/ariefzainuri96/ayo-test/cmd/api/dto/entity"

// @Model
type ProjectsResponse struct {
	BaseResponse
	Projects []entity.Project `json:"projects"`
	Pagination PaginationMetadata `json:"pagination"`
}

// @Model
type ProjectResponse struct {
	BaseResponse
	Project entity.Project `json:"project"`
}