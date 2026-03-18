package initiator

import (
	"context"
	"net/http"
	"os"
	"school-exam/internal/config"
	"school-exam/internal/db"
	"school-exam/internal/module/auth"
	"school-exam/internal/module/school"
	"school-exam/internal/module/superadmin"
	"school-exam/internal/module/teacher"
	"school-exam/internal/module/exam"
	"school-exam/internal/module/exam_session"
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var ct int
	if err := sqlDB.QueryRowContext(ctx, "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='users'").Scan(&ct); err != nil || ct == 0 {
		script, err := os.ReadFile("schemasqlite.sql")
		if err != nil {
			return nil, err
		}
		if err := db.ExecBatch(ctx, sqlDB, string(script)); err != nil {
			return nil, err
		}
	}
	// lightweight migration: ensure students.year exists
	var colCount int
	if err := sqlDB.QueryRowContext(ctx, "SELECT COUNT(*) FROM pragma_table_info('students') WHERE name = 'year'").Scan(&colCount); err == nil && colCount == 0 {
		_, _ = sqlDB.ExecContext(ctx, "ALTER TABLE students ADD COLUMN year TEXT")
	}
	if err := sqlDB.QueryRowContext(ctx, "SELECT COUNT(*) FROM pragma_table_info('exams') WHERE name = 'shuffle_options'").Scan(&colCount); err == nil && colCount == 0 {
		_, _ = sqlDB.ExecContext(ctx, "ALTER TABLE exams ADD COLUMN shuffle_options INTEGER DEFAULT 0")
	}

	_, _ = sqlDB.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS exam_sessions (
		id TEXT PRIMARY KEY,
		exam_id TEXT,
		student_id TEXT,
		start_time DATETIME,
		end_time DATETIME,
		status TEXT DEFAULT 'in_progress',
		total_score INTEGER,
		FOREIGN KEY (exam_id) REFERENCES exams(id),
		FOREIGN KEY (student_id) REFERENCES students(id)
	);`)

	_, _ = sqlDB.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS student_answers (
		id TEXT PRIMARY KEY,
		session_id TEXT,
		question_id TEXT,
		answer_text TEXT,
		selected_option_id TEXT,
		score INTEGER,
		FOREIGN KEY (session_id) REFERENCES exam_sessions(id),
		FOREIGN KEY (question_id) REFERENCES questions(id),
		FOREIGN KEY (selected_option_id) REFERENCES question_options(id),
		UNIQUE(session_id, question_id)
	);`)

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

	ts := security.TokenService{Secret: cfg.JWTSecret, TTL: time.Hour * 8}
	authUC := auth.NewAuthUsecase(usersRepo, ts)
	superUC := superadmin.NewUsecase(tenRepo, usersRepo, stuRepo, secRepo)
	schoolUC := school.NewUsecase(sqlDB, depRepo, secRepo, subRepo, teaRepo, usersRepo)
	teacherUC := teacher.NewUsecase(sqlDB, qbRepo, qqRepo, opRepo, teaRepo)
	examUC := exam.NewUsecase(sqlDB, exRepo, eqRepo, teaRepo, stuRepo, opRepo)
	sessionUC := exam_session.NewUsecase(sqlDB, esRepo, saRepo, stuRepo, exRepo, qqRepo, opRepo, vRepo)

	if err := authUC.SeedSuperAdmin(ctx, "superadmin", envDefault("SEED_SUPERADMIN_PASSWORD", "superadmin123")); err != nil {
		return nil, err
	}
	engine := route.SetupRouter(authUC, superUC, schoolUC, teacherUC, examUC, sessionUC, ts)
	s := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      engine,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
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
