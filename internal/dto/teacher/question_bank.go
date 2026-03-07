package teacherdto

type CreateQuestionBankRequest struct {
	SubjectID string `json:"subject_id" binding:"required"`
	Title     string `json:"title" binding:"required"`
}

type ListQuery struct {
	Page     int64 `form:"page"`
	PageSize int64 `form:"page_size"`
}

