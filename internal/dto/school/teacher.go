package schooldto

type CreateTeacherRequest struct {
	FirstName    string  `json:"first_name" binding:"required"`
	LastName     string  `json:"last_name" binding:"required"`
	TeacherCode  string  `json:"teacher_code" binding:"required"`
	DepartmentID *string `json:"department_id"`
	Email        *string `json:"email"`
}

type CreateTeacherResponse struct {
	TeacherID string `json:"teacher_id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type UpdateTeacherRequest struct {
	ID           string  `json:"id" binding:"required"`
	FirstName    *string `json:"first_name"`
	LastName     *string `json:"last_name"`
	TeacherCode  *string `json:"teacher_code"`
	DepartmentID *string `json:"department_id"`
}

type AssignRequest struct {
	TeacherID string `json:"teacher_id" binding:"required"`
	SubjectID string `json:"subject_id" binding:"required"`
	SectionID string `json:"section_id" binding:"required"`
}

type UnassignRequest struct {
	TeacherID string `json:"teacher_id" binding:"required"`
	SubjectID string `json:"subject_id" binding:"required"`
	SectionID string `json:"section_id" binding:"required"`
}

type AssignedTeacher struct {
	TeacherID    string  `json:"teacher_id"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	DepartmentID *string `json:"department_id"`
	SubjectID    string  `json:"subject_id"`
	SectionID    string  `json:"section_id"`
}

type TeacherListQuery struct {
	Page     int64 `form:"page"`
	PageSize int64 `form:"page_size"`
}
