package schooldto

type CreateSubjectRequest struct {
	Name         string  `json:"name" binding:"required"`
	DepartmentID *string `json:"department_id"`
}

type UpdateSubjectRequest struct {
	ID           string  `json:"id" binding:"required"`
	Name         *string `json:"name"`
	DepartmentID *string `json:"department_id"`
}

type SubjectListQuery struct {
	Page     int64 `form:"page"`
	PageSize int64 `form:"page_size"`
}
