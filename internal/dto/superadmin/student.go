package superadmindto

type CreateStudentRequest struct {
	TenantID     string  `json:"tenant_id" binding:"required"`
	StudentCode  string  `json:"student_code" binding:"required"`
	FirstName    string  `json:"first_name" binding:"required"`
	LastName     string  `json:"last_name" binding:"required"`
	Year         *string `json:"year"`
	SectionID    *string `json:"section_id"`
	DepartmentID *string `json:"department_id"`
	Email        *string `json:"email"`
}

type CreateStudentResponse struct {
	StudentID string `json:"student_id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type UpdateStudentRequest struct {
	ID           string  `json:"id" binding:"required"`
	FirstName    *string `json:"first_name"`
	LastName     *string `json:"last_name"`
	Year         *string `json:"year"`
	SectionID    *string `json:"section_id"`
	DepartmentID *string `json:"department_id"`
	Status       *string `json:"status"`
}

type ListQuery struct {
	Page     int64 `form:"page"`
	PageSize int64 `form:"page_size"`
}
