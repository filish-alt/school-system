package schooldto

type CreateDepartmentRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateDepartmentRequest struct {
	ID   string `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type ListQuery struct {
	Page     int64 `form:"page"`
	PageSize int64 `form:"page_size"`
}

