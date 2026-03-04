package schooldto

type CreateSectionRequest struct {
	Name         string  `json:"name" binding:"required"`
	DepartmentID *string `json:"department_id"`
	GradeLevel   *string `json:"grade_level"`
	AcademicYear *string `json:"academic_year"`
}

type UpdateSectionRequest struct {
	ID           string  `json:"id" binding:"required"`
	Name         *string `json:"name"`
	DepartmentID *string `json:"department_id"`
	GradeLevel   *string `json:"grade_level"`
	AcademicYear *string `json:"academic_year"`
}

type SectionListQuery struct {
	Page     int64 `form:"page"`
	PageSize int64 `form:"page_size"`
}
