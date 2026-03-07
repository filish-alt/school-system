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

	ts := security.TokenService{Secret: cfg.JWTSecret, TTL: time.Hour * 8}
	authUC := auth.NewAuthUsecase(usersRepo, ts)
	superUC := superadmin.NewUsecase(tenRepo, usersRepo, stuRepo)
	schoolUC := school.NewUsecase(sqlDB, depRepo, secRepo, subRepo, teaRepo, usersRepo)
	teacherUC := teacher.NewUsecase(sqlDB, qbRepo, qqRepo, opRepo, teaRepo)
	examUC := exam.NewUsecase(sqlDB, exRepo, eqRepo, teaRepo, opRepo)

	if err := authUC.SeedSuperAdmin(ctx, "superadmin", envDefault("SEED_SUPERADMIN_PASSWORD", "superadmin123")); err != nil {
		return nil, err
	}
	engine := route.SetupRouter(authUC, superUC, schoolUC, teacherUC, examUC, ts)
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
