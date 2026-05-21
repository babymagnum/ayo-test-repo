package utils

import (
	"fmt"
	"strings"

	"math"

	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/request"
	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/response"
	"gorm.io/gorm"
)

type PaginateResult[T any] struct {
	Data       []T
	Pagination response.PaginationMetadata
	Error      error
}

func ApplyPagination[T any](db *gorm.DB, req request.PaginationRequest, searchAllQuery string) PaginateResult[T] {

	countQuery := db.Session(&gorm.Session{})

	paginatedQuery := db.Session(&gorm.Session{})

	if req.SearchField != "" && req.SearchValue != "" {
		whereClause := fmt.Sprintf("CAST(%s AS TEXT) ILIKE ?", req.SearchField)

		countQuery = countQuery.Where(whereClause, "%"+req.SearchValue+"%")
		paginatedQuery = paginatedQuery.Where(whereClause, "%"+req.SearchValue+"%")
	}

	if req.SearchAll != "" && searchAllQuery != "" {
		search := "%" + req.SearchAll + "%"
		howManyFields := strings.Count(searchAllQuery, "?")

		if howManyFields > 0 {
			args := make([]interface{}, howManyFields)
			for i := range howManyFields {
				args[i] = search
			}

			countQuery = countQuery.Where(searchAllQuery, args...)
			paginatedQuery = paginatedQuery.Where(searchAllQuery, args...)
		}
	}

	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return PaginateResult[T]{Error: fmt.Errorf("failed to count records: %w", err)}
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))
	offset := (req.Page - 1) * req.PageSize

	orderBy := "id"
	if req.OrderBy != "" {
		orderBy = req.OrderBy
	}

	sortDirection := "ASC"
	if strings.ToUpper(req.Sort) == "DESC" {
		sortDirection = "DESC"
	}

	paginatedQuery = paginatedQuery.Order(fmt.Sprintf("%s %s", orderBy, sortDirection))

	paginatedQuery = paginatedQuery.Offset(offset).Limit(req.PageSize)

	var data []T
	if err := paginatedQuery.Find(&data).Error; err != nil {
		return PaginateResult[T]{Error: fmt.Errorf("failed to fetch records: %w", err)}
	}

	metadata := response.PaginationMetadata{
		Page:      req.Page,
		PageSize:  req.PageSize,
		TotalData: total,
		TotalPage: totalPages,
	}

	return PaginateResult[T]{
		Data:       data,
		Pagination: metadata,
	}
}
