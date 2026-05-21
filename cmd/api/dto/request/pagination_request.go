package request

type PaginationRequest struct {
	Page        int    `form:"page" binding:"required"`
	PageSize    int    `form:"page_size" binding:"required"`
	SearchField string `form:"search_field"`
	SearchValue string `form:"search_value"`
	SearchAll   string `form:"search_all"`
	OrderBy     string `form:"order_by"`
	Sort        string `form:"sort"`
}
