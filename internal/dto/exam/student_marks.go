package examdto

type StudentMark struct {
	StudentCode string `json:"student_code"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	SectionName string `json:"section_name"`
	TotalScore  int64  `json:"total_score"`
}

type ExamMarksResponse struct {
	ExamTitle string        `json:"exam_title"`
	Marks     []StudentMark `json:"marks"`
}
