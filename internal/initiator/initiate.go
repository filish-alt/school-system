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
	usersRepo := repository.NewUserRepository(sqlDB)
	tenRepo := repository.NewTenantRepository(sqlDB)
	stuRepo := repository.NewStudentRepository(sqlDB)
	depRepo := repository.NewDepartmentRepository(sqlDB)
	secRepo := repository.NewSectionRepository(sqlDB)
	subRepo := repository.NewSubjectRepository(sqlDB)
	teaRepo := repository.NewTeacherRepository(sqlDB)
	ts := security.TokenService{Secret: cfg.JWTSecret, TTL: time.Hour * 8}
	authUC := auth.NewAuthUsecase(usersRepo, ts)
	superUC := superadmin.NewUsecase(tenRepo, usersRepo, stuRepo)
	schoolUC := school.NewUsecase(sqlDB, depRepo, secRepo, subRepo, teaRepo, usersRepo)
	if err := authUC.SeedSuperAdmin(ctx, "superadmin", envDefault("SEED_SUPERADMIN_PASSWORD", "superadmin123")); err != nil {
		return nil, err
	}
	engine := route.SetupRouter(authUC, superUC, schoolUC, ts)
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
