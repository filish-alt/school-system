package system

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"school-exam/internal/domain"
	"school-exam/internal/repository"
	"time"

	"github.com/google/uuid"
)

type Usecase struct {
	db        *sql.DB
	auditRepo *repository.AuditLogRepository
	exRepo    *repository.ExamRepository
	eqRepo    *repository.ExamQuestionRepository
	stuRepo   *repository.StudentRepository
}

func NewUsecase(db *sql.DB, auditRepo *repository.AuditLogRepository, exRepo *repository.ExamRepository, eqRepo *repository.ExamQuestionRepository, stuRepo *repository.StudentRepository) *Usecase {
	return &Usecase{
		db:        db,
		auditRepo: auditRepo,
		exRepo:    exRepo,
		eqRepo:    eqRepo,
		stuRepo:   stuRepo,
	}
}

func (u *Usecase) LogAction(ctx context.Context, tenantID *string, userID *string, action string, entityType, entityID, details, ip *string) {
	log := domain.AuditLog{
		ID:         uuid.NewString(),
		TenantID:   tenantID,
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Details:    details,
		IPAddress:  ip,
	}
	_ = u.auditRepo.Log(ctx, log)
}

func (u *Usecase) ListAuditLogs(ctx context.Context, tenantID string, page, pageSize int) ([]domain.AuditLog, error) {
	if page < 1 { page = 1 }
	if pageSize < 1 { pageSize = 50 }
	offset := (page - 1) * pageSize
	return u.auditRepo.List(ctx, tenantID, pageSize, offset)
}

// Backup & Recovery
type BackupInfo struct {
	Filename  string    `json:"filename"`
	SizeBytes int64     `json:"size_bytes"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *Usecase) CreateBackup(ctx context.Context, dbPath string) (string, error) {
	_ = os.MkdirAll("backups", 0755)
	timestamp := time.Now().Format("20060102_150405")
	backupFile := filepath.Join("backups", fmt.Sprintf("school_backup_%s.db", timestamp))

	srcFile, err := os.Open(dbPath)
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	destFile, err := os.Create(backupFile)
	if err != nil {
		return "", err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return "", err
	}

	return backupFile, nil
}

func (u *Usecase) ListBackups() ([]BackupInfo, error) {
	_ = os.MkdirAll("backups", 0755)
	files, err := os.ReadDir("backups")
	if err != nil {
		return nil, err
	}

	var list []BackupInfo
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".db" {
			info, err := f.Info()
			if err == nil {
				list = append(list, BackupInfo{
					Filename:  f.Name(),
					SizeBytes: info.Size(),
					CreatedAt: info.ModTime(),
				})
			}
		}
	}
	return list, nil
}

// Emergency Exam Package & Paper Fallback Data
type EmergencyExamPackage struct {
	ExamTitle        string        `json:"exam_title"`
	SubjectName      string        `json:"subject_name"`
	SectionName      string        `json:"section_name"`
	DurationMinutes  int64         `json:"duration_minutes"`
	TotalMarks       float64       `json:"total_marks"`
	Instructions     string        `json:"instructions"`
	AcademicYear     string        `json:"academic_year"`
	GeneratedAt      time.Time     `json:"generated_at"`
	Questions        []PaperQuestion `json:"questions"`
	StudentRoster    []PaperStudent  `json:"student_roster"`
}

type PaperQuestion struct {
	QuestionNumber int      `json:"question_number"`
	Type           string   `json:"type"`
	QuestionText   string   `json:"question_text"`
	Marks          float64  `json:"marks"`
	Options        []string `json:"options,omitempty"`
}

type PaperStudent struct {
	StudentID   string `json:"student_id"`
	StudentCode string `json:"student_code"`
	FullName    string `json:"full_name"`
}

func (u *Usecase) GenerateEmergencyPackage(ctx context.Context, examID string) (*EmergencyExamPackage, error) {
	// 1. Get Exam details
	queryExam := `SELECT e.title, sub.name as subject_name, sec.name as section_name, e.duration_minutes,
	                     COALESCE(e.total_marks, 100), COALESCE(e.instructions, ''), COALESCE(sec.academic_year, '2025/2026'), e.section_id
	              FROM exams e
	              JOIN subjects sub ON e.subject_id = sub.id
	              JOIN sections sec ON e.section_id = sec.id
	              WHERE e.id = ?`
	var pkg EmergencyExamPackage
	var secID string
	err := u.db.QueryRowContext(ctx, queryExam, examID).Scan(
		&pkg.ExamTitle, &pkg.SubjectName, &pkg.SectionName, &pkg.DurationMinutes,
		&pkg.TotalMarks, &pkg.Instructions, &pkg.AcademicYear, &secID,
	)
	if err != nil {
		return nil, err
	}
	pkg.GeneratedAt = time.Now()

	// 2. Get Exam Questions
	queryQuestions := `SELECT q.id, q.type, q.question_text, eq.marks, eq.order_index
	                   FROM exam_questions eq
	                   JOIN questions q ON eq.question_id = q.id
	                   WHERE eq.exam_id = ? AND eq.is_dropped = 0
	                   ORDER BY eq.order_index ASC`
	rows, err := u.db.QueryContext(ctx, queryQuestions, examID)
	if err == nil {
		defer rows.Close()
		idx := 1
		for rows.Next() {
			var qID, qType, qText string
			var qMarks float64
			var orderIdx sql.NullInt64
			if err := rows.Scan(&qID, &qType, &qText, &qMarks, &orderIdx); err == nil {
				pq := PaperQuestion{
					QuestionNumber: idx,
					Type:           qType,
					QuestionText:   qText,
					Marks:          qMarks,
				}
				// Load options for MCQ
				optRows, err := u.db.QueryContext(ctx, "SELECT option_text FROM question_options WHERE question_id = ? ORDER BY order_index ASC", qID)
				if err == nil {
					for optRows.Next() {
						var optText string
						if err := optRows.Scan(&optText); err == nil {
							pq.Options = append(pq.Options, optText)
						}
					}
					optRows.Close()
				}
				pkg.Questions = append(pkg.Questions, pq)
				idx++
			}
		}
	}

	// 3. Get Student Roster
	queryStudents := `SELECT id, student_code, (first_name || ' ' || last_name) as name FROM students WHERE section_id = ? ORDER BY student_code ASC`
	stuRows, err := u.db.QueryContext(ctx, queryStudents, secID)
	if err == nil {
		defer stuRows.Close()
		for stuRows.Next() {
			var s PaperStudent
			if err := stuRows.Scan(&s.StudentID, &s.StudentCode, &s.FullName); err == nil {
				pkg.StudentRoster = append(pkg.StudentRoster, s)
			}
		}
	}

	return &pkg, nil
}
