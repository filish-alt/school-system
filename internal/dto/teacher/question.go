package teacherdto

type CreateQuestionRequest struct {
	QuestionBankID string  `json:"question_bank_id" binding:"required"`
	Type           string  `json:"type" binding:"required"`
	QuestionText   string  `json:"question_text" binding:"required"`
	ImageURL       *string `json:"image_url"`
	Marks          *int64  `json:"marks"`
	Difficulty     *string `json:"difficulty_level"`
}

type UpdateQuestionRequest struct {
	ID             string  `json:"id" binding:"required"`
	Type           *string `json:"type"`
	QuestionText   *string `json:"question_text"`
	ImageURL       *string `json:"image_url"`
	Marks          *int64  `json:"marks"`
	Difficulty     *string `json:"difficulty_level"`
}

type CreateOptionRequest struct {
	QuestionID string `json:"question_id" binding:"required"`
	OptionText string `json:"option_text" binding:"required"`
	IsCorrect  bool   `json:"is_correct"`
}

type UpdateOptionRequest struct {
	ID         string  `json:"id" binding:"required"`
	OptionText *string `json:"option_text"`
	IsCorrect  *bool   `json:"is_correct"`
}

