package domain

import "time"

type Role struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID                 string     `json:"id"`
	TenantID           *string    `json:"tenant_id,omitempty"`
	Username           string     `json:"username"`
	PasswordHash       string     `json:"-"`
	Email              *string    `json:"email,omitempty"`
	Phone              *string    `json:"phone,omitempty"`
	ProfilePhoto       *string    `json:"profile_photo,omitempty"`
	RoleID             *int64     `json:"role_id,omitempty"`
	RoleName           *string    `json:"role_name,omitempty"`
	Status             string     `json:"status"`
	MustChangePassword bool       `json:"must_change_password"`
	CreatedAt          time.Time  `json:"created_at"`
}

type AcademicYear struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	YearName  string    `json:"year_name"`
	StartDate *string   `json:"start_date,omitempty"`
	EndDate   *string   `json:"end_date,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type SystemSetting struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenant_id"`
	SettingKey   string    `json:"setting_key"`
	SettingValue string    `json:"setting_value"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type TeacherSchedule struct {
	ID           string  `json:"id"`
	TenantID     string  `json:"tenant_id"`
	TeacherID    string  `json:"teacher_id"`
	TeacherName  *string `json:"teacher_name,omitempty"`
	SubjectID    *string `json:"subject_id,omitempty"`
	SubjectName  *string `json:"subject_name,omitempty"`
	SectionID    *string `json:"section_id,omitempty"`
	SectionName  *string `json:"section_name,omitempty"`
	DayOfWeek    string  `json:"day_of_week"`
	StartTime    string  `json:"start_time"`
	EndTime      string  `json:"end_time"`
	Room         *string `json:"room,omitempty"`
	ActivityType string  `json:"activity_type"`
}

type CourseworkAssessment struct {
	ID               string    `json:"id"`
	TenantID         string    `json:"tenant_id"`
	AcademicYearID   *string   `json:"academic_year_id,omitempty"`
	SubjectID        string    `json:"subject_id"`
	SubjectName      *string   `json:"subject_name,omitempty"`
	SectionID        string    `json:"section_id"`
	SectionName      *string   `json:"section_name,omitempty"`
	TeacherID        string    `json:"teacher_id"`
	TeacherName      *string   `json:"teacher_name,omitempty"`
	Title            string    `json:"title"`
	Type             string    `json:"type"` // quiz, mid_exam, assignment, project
	MaxMarks         float64   `json:"max_marks"`
	WeightPercentage float64   `json:"weight_percentage"`
	DueDate          *string   `json:"due_date,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

type StudentCourseworkMark struct {
	ID            string    `json:"id"`
	CourseworkID  string    `json:"coursework_id"`
	StudentID     string    `json:"student_id"`
	StudentCode   *string   `json:"student_code,omitempty"`
	StudentName   *string   `json:"student_name,omitempty"`
	MarksObtained float64   `json:"marks_obtained"`
	Remarks       *string   `json:"remarks,omitempty"`
	GradedBy      *string   `json:"graded_by,omitempty"`
	GradedAt      time.Time `json:"graded_at"`
}

type ExamExtension struct {
	ID              string    `json:"id"`
	ExamID          string    `json:"exam_id"`
	StudentID       *string   `json:"student_id,omitempty"`
	StudentName     *string   `json:"student_name,omitempty"`
	ExtendedMinutes int       `json:"extended_minutes"`
	Reason          string    `json:"reason"`
	AuthorizedBy    string    `json:"authorized_by"`
	AuthorizerName  *string   `json:"authorizer_name,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type ExamRemoteAuthorization struct {
	ID             string    `json:"id"`
	ExamID         string    `json:"exam_id"`
	StudentID      string    `json:"student_id"`
	StudentCode    *string   `json:"student_code,omitempty"`
	StudentName    *string   `json:"student_name,omitempty"`
	ValidUntil     string    `json:"valid_until"`
	Reason         string    `json:"reason"`
	ApprovedBy     string    `json:"approved_by"`
	AuthorizerName *string   `json:"authorizer_name,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type DroppedQuestion struct {
	ID           string    `json:"id"`
	ExamID       string    `json:"exam_id"`
	QuestionID   string    `json:"question_id"`
	QuestionText *string   `json:"question_text,omitempty"`
	Reason       string    `json:"reason"`
	CreditType   string    `json:"credit_type"` // full_credit, reduce_total
	DroppedBy    string    `json:"dropped_by"`
	CreatedAt    time.Time `json:"created_at"`
}

type AnswerKeyCorrection struct {
	ID           string    `json:"id"`
	ExamID       string    `json:"exam_id"`
	QuestionID   string    `json:"question_id"`
	QuestionText *string   `json:"question_text,omitempty"`
	OldAnswer    *string   `json:"old_answer,omitempty"`
	NewAnswer    *string   `json:"new_answer,omitempty"`
	Reason       string    `json:"reason"`
	CorrectedBy  string    `json:"corrected_by"`
	CreatedAt    time.Time `json:"created_at"`
}

type StudentComplaint struct {
	ID              string     `json:"id"`
	TenantID        string     `json:"tenant_id"`
	StudentID       string     `json:"student_id"`
	StudentCode     *string    `json:"student_code,omitempty"`
	StudentName     *string    `json:"student_name,omitempty"`
	SubjectID       *string    `json:"subject_id,omitempty"`
	SubjectName     *string    `json:"subject_name,omitempty"`
	ExamID          *string    `json:"exam_id,omitempty"`
	ExamTitle       *string    `json:"exam_title,omitempty"`
	CourseworkID    *string    `json:"coursework_id,omitempty"`
	QuestionID      *string    `json:"question_id,omitempty"`
	Category        string     `json:"category"`
	Description     string     `json:"description"`
	Status          string     `json:"status"` // pending, escalated_to_dept_head, resolved, rejected
	Response        *string    `json:"response,omitempty"`
	ResolvedBy      *string    `json:"resolved_by,omitempty"`
	ResolverName    *string    `json:"resolver_name,omitempty"`
	ResolvedAt      *time.Time `json:"resolved_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type QuestionMatchPair struct {
	ID          string `json:"id"`
	QuestionID  string `json:"question_id"`
	LeftPrompt  string `json:"left_prompt"`
	RightMatch  string `json:"right_match"`
	OrderIndex  int    `json:"order_index"`
}

type AssessmentResult struct {
	ID              string     `json:"id"`
	TenantID        *string    `json:"tenant_id,omitempty"`
	AcademicYearID  *string    `json:"academic_year_id,omitempty"`
	StudentID       string     `json:"student_id"`
	StudentCode     string     `json:"student_code"`
	StudentName     string     `json:"student_name"`
	SubjectID       string     `json:"subject_id"`
	SubjectName     string     `json:"subject_name"`
	SectionID       *string    `json:"section_id,omitempty"`
	SectionName     *string    `json:"section_name,omitempty"`
	QuizScore       float64    `json:"quiz_score"`
	MidScore        float64    `json:"mid_score"`
	AssignmentScore float64    `json:"assignment_score"`
	ProjectScore    float64    `json:"project_score"`
	CAScore         float64    `json:"ca_score"`
	FinalScore      float64    `json:"final_score"`
	TotalScore      float64    `json:"total_score"`
	Grade           string     `json:"grade"`
	IsPublished     bool       `json:"is_published"`
	PublishedAt     *time.Time `json:"published_at,omitempty"`
}

type AuditLog struct {
	ID         string    `json:"id"`
	TenantID   *string   `json:"tenant_id,omitempty"`
	UserID     *string   `json:"user_id,omitempty"`
	Username   *string   `json:"username,omitempty"`
	Action     string    `json:"action"`
	EntityType *string   `json:"entity_type,omitempty"`
	EntityID   *string   `json:"entity_id,omitempty"`
	Details    *string   `json:"details,omitempty"`
	IPAddress  *string   `json:"ip_address,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
