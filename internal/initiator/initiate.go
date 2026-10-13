package initiator

import (
	"context"
	"net/http"
	"os"
	"school-exam/internal/config"
	"school-exam/internal/db"
	"school-exam/internal/module/academic"
	"school-exam/internal/module/auth"
	"school-exam/internal/module/complaint"
	"school-exam/internal/module/coursework"
	"school-exam/internal/module/exam"
	"school-exam/internal/module/exam_control"
	"school-exam/internal/module/exam_session"
	"school-exam/internal/module/result"
	"school-exam/internal/module/schedule"
	"school-exam/internal/module/school"
	"school-exam/internal/module/superadmin"
	"school-exam/internal/module/system"
	"school-exam/internal/module/teacher"
	"school-exam/internal/repository"
	"school-exam/internal/route"
	"school-exam/internal/security"
	"time"
)

type App struct {
	Server *http.Server
}

func Initiate() (*App, error) {
	cfg := config.Load()
	sqlDB, err := db.Open(cfg.DBPath)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// 1. Run full schema script
	script, err := os.ReadFile("schemasqlite.sql")
	if err == nil {
		_ = db.ExecBatch(ctx, sqlDB, string(script))
	}

	// 2. Safe migrations for existing databases
	migrations := []string{
		"ALTER TABLE students ADD COLUMN year TEXT",
		"ALTER TABLE students ADD COLUMN phone TEXT",
		"ALTER TABLE students ADD COLUMN photo_url TEXT",
		"ALTER TABLE students ADD COLUMN academic_year_id TEXT",
		"ALTER TABLE users ADD COLUMN phone TEXT",
		"ALTER TABLE users ADD COLUMN profile_photo TEXT",
		"ALTER TABLE users ADD COLUMN must_change_password INTEGER DEFAULT 0",
		"ALTER TABLE teachers ADD COLUMN phone TEXT",
		"ALTER TABLE teachers ADD COLUMN photo_url TEXT",
		"ALTER TABLE sections ADD COLUMN academic_year_id TEXT",
		"ALTER TABLE exams ADD COLUMN exam_type TEXT DEFAULT 'mid'",
		"ALTER TABLE exams ADD COLUMN approval_status TEXT DEFAULT 'draft'",
		"ALTER TABLE exams ADD COLUMN is_locked INTEGER DEFAULT 0",
		"ALTER TABLE exams ADD COLUMN is_visible INTEGER DEFAULT 0",
		"ALTER TABLE exams ADD COLUMN is_started INTEGER DEFAULT 0",
		"ALTER TABLE exams ADD COLUMN access_code TEXT",
		"ALTER TABLE exams ADD COLUMN instructions TEXT",
		"ALTER TABLE exams ADD COLUMN pass_percentage REAL DEFAULT 50.0",
		"ALTER TABLE exams ADD COLUMN shuffle_options INTEGER DEFAULT 0",
		"ALTER TABLE exams ADD COLUMN version INTEGER DEFAULT 1",
		"ALTER TABLE exams ADD COLUMN approved_by TEXT",
		"ALTER TABLE exams ADD COLUMN approved_at DATETIME",
		"ALTER TABLE exams ADD COLUMN academic_year_id TEXT",
		"ALTER TABLE exams ADD COLUMN department_id TEXT",
		"ALTER TABLE questions ADD COLUMN topic TEXT",
		"ALTER TABLE questions ADD COLUMN grade_level TEXT",
		"ALTER TABLE questions ADD COLUMN academic_year_id TEXT",
		"ALTER TABLE questions ADD COLUMN creator_teacher_id TEXT",
		"ALTER TABLE questions ADD COLUMN is_private INTEGER DEFAULT 1",
		"ALTER TABLE questions ADD COLUMN status TEXT DEFAULT 'draft'",
		"ALTER TABLE questions ADD COLUMN answer_key TEXT",
		"ALTER TABLE questions ADD COLUMN case_sensitive INTEGER DEFAULT 0",
		"ALTER TABLE exam_questions ADD COLUMN is_dropped INTEGER DEFAULT 0",
		"ALTER TABLE exam_sessions ADD COLUMN is_remote INTEGER DEFAULT 0",
		"ALTER TABLE student_answers ADD COLUMN matching_answers TEXT",
		"ALTER TABLE student_answers ADD COLUMN is_manually_graded INTEGER DEFAULT 0",
		`CREATE TABLE IF NOT EXISTS audit_logs (
			id TEXT PRIMARY KEY,
			tenant_id TEXT,
			user_id TEXT,
			action TEXT NOT NULL,
			entity_type TEXT,
			entity_id TEXT,
			details TEXT,
			ip_address TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS results (
			id TEXT PRIMARY KEY,
			tenant_id TEXT,
			academic_year_id TEXT,
			student_id TEXT NOT NULL,
			subject_id TEXT NOT NULL,
			section_id TEXT,
			quiz_score REAL DEFAULT 0,
			mid_score REAL DEFAULT 0,
			assignment_score REAL DEFAULT 0,
			project_score REAL DEFAULT 0,
			ca_score REAL DEFAULT 0,
			final_score REAL DEFAULT 0,
			total_score REAL DEFAULT 0,
			grade TEXT,
			is_published INTEGER DEFAULT 0,
			published_at DATETIME,
			UNIQUE(student_id, subject_id, academic_year_id)
		)`,
		"ALTER TABLE results ADD COLUMN tenant_id TEXT",
		"ALTER TABLE results ADD COLUMN academic_year_id TEXT",
		"ALTER TABLE results ADD COLUMN student_id TEXT",
		"ALTER TABLE results ADD COLUMN subject_id TEXT",
		"ALTER TABLE results ADD COLUMN section_id TEXT",
		"ALTER TABLE results ADD COLUMN quiz_score REAL DEFAULT 0",
		"ALTER TABLE results ADD COLUMN mid_score REAL DEFAULT 0",
		"ALTER TABLE results ADD COLUMN assignment_score REAL DEFAULT 0",
		"ALTER TABLE results ADD COLUMN project_score REAL DEFAULT 0",
		"ALTER TABLE results ADD COLUMN ca_score REAL DEFAULT 0",
		"ALTER TABLE results ADD COLUMN final_score REAL DEFAULT 0",
		"ALTER TABLE results ADD COLUMN total_score REAL DEFAULT 0",
		"ALTER TABLE results ADD COLUMN grade TEXT",
		"ALTER TABLE results ADD COLUMN is_published INTEGER DEFAULT 0",
		"ALTER TABLE results ADD COLUMN published_at DATETIME",
		`CREATE TABLE IF NOT EXISTS system_settings (
			id TEXT PRIMARY KEY,
			tenant_id TEXT,
			setting_key TEXT NOT NULL,
			setting_value TEXT NOT NULL,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(tenant_id, setting_key)
		)`,
		"ALTER TABLE system_settings ADD COLUMN tenant_id TEXT",
		`CREATE TABLE IF NOT EXISTS academic_years (
			id TEXT PRIMARY KEY,
			tenant_id TEXT,
			name TEXT NOT NULL,
			start_date DATE NOT NULL,
			end_date DATE NOT NULL,
			is_active INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS coursework_assessments (
			id TEXT PRIMARY KEY,
			tenant_id TEXT,
			academic_year_id TEXT,
			subject_id TEXT NOT NULL,
			section_id TEXT NOT NULL,
			teacher_id TEXT NOT NULL,
			title TEXT NOT NULL,
			type TEXT NOT NULL,
			max_marks REAL NOT NULL,
			weight_percentage REAL NOT NULL,
			due_date DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS student_coursework_marks (
			id TEXT PRIMARY KEY,
			coursework_id TEXT NOT NULL,
			student_id TEXT NOT NULL,
			marks_obtained REAL NOT NULL,
			remarks TEXT,
			graded_by TEXT,
			graded_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(coursework_id, student_id)
		)`,
		`CREATE TABLE IF NOT EXISTS dropped_questions (
			id TEXT PRIMARY KEY,
			exam_id TEXT NOT NULL,
			question_id TEXT NOT NULL,
			reason TEXT NOT NULL,
			credit_type TEXT DEFAULT 'full_credit',
			dropped_by TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS answer_key_corrections (
			id TEXT PRIMARY KEY,
			exam_id TEXT NOT NULL,
			question_id TEXT NOT NULL,
			old_answer TEXT,
			new_answer TEXT,
			reason TEXT NOT NULL,
			corrected_by TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS student_complaints (
			id TEXT PRIMARY KEY,
			tenant_id TEXT NOT NULL,
			student_id TEXT NOT NULL,
			subject_id TEXT,
			exam_id TEXT,
			coursework_id TEXT,
			question_id TEXT,
			category TEXT,
			description TEXT NOT NULL,
			status TEXT DEFAULT 'pending',
			response TEXT,
			resolved_by TEXT,
			resolved_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS timetables (
			id TEXT PRIMARY KEY,
			tenant_id TEXT,
			teacher_id TEXT NOT NULL,
			subject_id TEXT,
			section_id TEXT,
			day_of_week TEXT NOT NULL,
			start_time TEXT NOT NULL,
			end_time TEXT NOT NULL,
			room TEXT,
			activity_type TEXT DEFAULT 'lecture',
			academic_year_id TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, m := range migrations {
		_, _ = sqlDB.ExecContext(ctx, m)
	}

	// 3. Ensure all standard roles exist
	roles := []string{
		"super_admin",
		"school_director",
		"vice_director",
		"department_head",
		"teacher",
		"student",
		"school_admin",
	}
	for _, r := range roles {
		_, _ = sqlDB.ExecContext(ctx, "INSERT OR IGNORE INTO roles (name) VALUES (?)", r)
	}

	// Repositories
	usersRepo := repository.NewUserRepository(sqlDB)
	tenRepo := repository.NewTenantRepository(sqlDB)
	stuRepo := repository.NewStudentRepository(sqlDB)
	depRepo := repository.NewDepartmentRepository(sqlDB)
	secRepo := repository.NewSectionRepository(sqlDB)
	subRepo := repository.NewSubjectRepository(sqlDB)
	teaRepo := repository.NewTeacherRepository(sqlDB)
	qbRepo := repository.NewQuestionBankRepository(sqlDB)
	qqRepo := repository.NewQuestionRepository(sqlDB)
	opRepo := repository.NewOptionRepository(sqlDB)
	exRepo := repository.NewExamRepository(sqlDB)
	eqRepo := repository.NewExamQuestionRepository(sqlDB)
	esRepo := repository.NewExamSessionRepository(sqlDB)
	saRepo := repository.NewStudentAnswerRepository(sqlDB)
	vRepo := repository.NewExamViolationRepository(sqlDB)

	ayRepo := repository.NewAcademicYearRepository(sqlDB)
	setRepo := repository.NewSystemSettingRepository(sqlDB)
	cwRepo := repository.NewCourseworkRepository(sqlDB)
	ecRepo := repository.NewExamControlRepository(sqlDB)
	compRepo := repository.NewComplaintRepository(sqlDB)
	schRepo := repository.NewScheduleRepository(sqlDB)
	resRepo := repository.NewResultRepository(sqlDB)
	auditRepo := repository.NewAuditLogRepository(sqlDB)

	// Usecases
	ts := security.TokenService{Secret: cfg.JWTSecret, TTL: time.Hour * 8}
	authUC := auth.NewAuthUsecase(usersRepo, ts)
	superUC := superadmin.NewUsecase(tenRepo, usersRepo, stuRepo, secRepo)
	schoolUC := school.NewUsecase(sqlDB, depRepo, secRepo, subRepo, teaRepo, usersRepo)
	teacherUC := teacher.NewUsecase(sqlDB, qbRepo, qqRepo, opRepo, teaRepo)
	examUC := exam.NewUsecase(sqlDB, exRepo, eqRepo, teaRepo, stuRepo, opRepo)
	sessionUC := exam_session.NewUsecase(sqlDB, esRepo, saRepo, stuRepo, exRepo, qqRepo, opRepo, vRepo)

	acadUC := academic.NewUsecase(ayRepo, setRepo)
	cwUC := coursework.NewUsecase(cwRepo)
	ecUC := exam_control.NewUsecase(ecRepo)
	compUC := complaint.NewUsecase(compRepo)
	schUC := schedule.NewUsecase(schRepo)
	resUC := result.NewUsecase(resRepo)
	sysUC := system.NewUsecase(sqlDB, auditRepo, exRepo, eqRepo, stuRepo)

	if err := authUC.SeedSuperAdmin(ctx, "superadmin", envDefault("SEED_SUPERADMIN_PASSWORD", "superadmin123")); err != nil {
		return nil, err
	}

	engine := route.SetupRouter(
		authUC, superUC, schoolUC, teacherUC, examUC, sessionUC,
		acadUC, cwUC, ecUC, compUC, schUC, resUC, sysUC,
		ts, cfg,
	)

	s := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return &App{Server: s}, nil
}

func envDefault(k, d string) string {
	v := os.Getenv(k)
	if v == "" {
		return d
	}
	return v
}
