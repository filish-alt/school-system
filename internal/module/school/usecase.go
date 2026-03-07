package school

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"strings"

	"school-exam/internal/domain"
	schooldto "school-exam/internal/dto/school"
	"school-exam/internal/repository"
	"school-exam/internal/security"
	q "school-exam/internal/sqlc/gen"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Usecase struct {
	Departments *repository.DepartmentRepository
	Sections    *repository.SectionRepository
	Subjects    *repository.SubjectRepository
	Teachers    *repository.TeacherRepository
	Users       *repository.UserRepository
	Queries     *q.Queries
}

func NewUsecase(db *sql.DB, deps *repository.DepartmentRepository, secs *repository.SectionRepository, subs *repository.SubjectRepository, teachers *repository.TeacherRepository, users *repository.UserRepository) *Usecase {
	return &Usecase{
		Departments: deps,
		Sections:    secs,
		Subjects:    subs,
		Teachers:    teachers,
		Users:       users,
		Queries:     q.New(db),
	}
}

func (u *Usecase) tenantFromContext(ctx context.Context) (string, error) {
	c := ctx.(*gin.Context)
	v, ok := c.Get("claims")
	if !ok {
		return "", fmt.Errorf("claims missing")
	}
	claims := v.(*security.Claims)
	if claims.TenantID == nil || *claims.TenantID == "" {
		return "", fmt.Errorf("tenant missing")
	}
	return *claims.TenantID, nil
}

// Departments
func (u *Usecase) CreateDepartment(ctx context.Context, req schooldto.CreateDepartmentRequest) (string, error) {
	tid, err := u.tenantFromContext(ctx)
	if err != nil {
		return "", err
	}
	id := uuid.New().String()
	if err := u.Departments.Create(ctx, id, tid, req.Name); err != nil {
		return "", err
	}
	return id, nil
}

func (u *Usecase) ListDepartments(ctx context.Context, page, pageSize int64) ([]q.Department, error) {
	tid, err := u.tenantFromContext(ctx)
	if err != nil {
		return nil, err
	}
	limit, offset := normalize(page, pageSize)
	return u.Departments.ListByTenant(ctx, tid, limit, offset)
}

func (u *Usecase) UpdateDepartment(ctx context.Context, req schooldto.UpdateDepartmentRequest) error {
	return u.Departments.Update(ctx, req.ID, req.Name)
}

func (u *Usecase) DeleteDepartment(ctx context.Context, id string) error {
	return u.Departments.Delete(ctx, id)
}

// Sections
func (u *Usecase) CreateSection(ctx context.Context, req schooldto.CreateSectionRequest) (string, error) {
	tid, err := u.tenantFromContext(ctx)
	if err != nil {
		return "", err
	}
	id := uuid.New().String()
	var dep, grade, year sql.NullString
	if req.DepartmentID != nil {
		dep = sql.NullString{String: *req.DepartmentID, Valid: true}
	}
	if req.GradeLevel != nil {
		grade = sql.NullString{String: *req.GradeLevel, Valid: true}
	}
	if req.AcademicYear != nil {
		year = sql.NullString{String: *req.AcademicYear, Valid: true}
	}
	if err := u.Sections.Create(ctx, q.CreateSectionParams{
		ID:           id,
		TenantID:     sql.NullString{String: tid, Valid: true},
		Name:         req.Name,
		DepartmentID: dep,
		GradeLevel:   grade,
		AcademicYear: year,
	}); err != nil {
		return "", err
	}
	return id, nil
}

func (u *Usecase) ListSections(ctx context.Context, page, pageSize int64) ([]q.Section, error) {
	tid, err := u.tenantFromContext(ctx)
	if err != nil {
		return nil, err
	}
	limit, offset := normalize(page, pageSize)
	return u.Sections.ListByTenant(ctx, tid, limit, offset)
}

func (u *Usecase) UpdateSection(ctx context.Context, req schooldto.UpdateSectionRequest) error {
	// fetch current to obtain existing name when not provided
	cur, err := u.Sections.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}
	newName := cur.Name
	if req.Name != nil {
		newName = *req.Name
	}
	var dep, grade, year sql.NullString
	if req.DepartmentID != nil {
		dep = sql.NullString{String: *req.DepartmentID, Valid: true}
	}
	if req.GradeLevel != nil {
		grade = sql.NullString{String: *req.GradeLevel, Valid: true}
	}
	if req.AcademicYear != nil {
		year = sql.NullString{String: *req.AcademicYear, Valid: true}
	}
	return u.Sections.Update(ctx, q.UpdateSectionParams{
		ID:           req.ID,
		Name:         newName,
		DepartmentID: dep,
		GradeLevel:   grade,
		AcademicYear: year,
	})
}

func (u *Usecase) DeleteSection(ctx context.Context, id string) error {
	return u.Sections.Delete(ctx, id)
}

// Subjects
func (u *Usecase) CreateSubject(ctx context.Context, req schooldto.CreateSubjectRequest) (string, error) {
	tid, err := u.tenantFromContext(ctx)
	if err != nil {
		return "", err
	}
	id := uuid.New().String()
	var dep sql.NullString
	if req.DepartmentID != nil {
		dep = sql.NullString{String: *req.DepartmentID, Valid: true}
	}
	if err := u.Subjects.Create(ctx, q.CreateSubjectParams{
		ID:           id,
		TenantID:     sql.NullString{String: tid, Valid: true},
		Name:         req.Name,
		DepartmentID: dep,
	}); err != nil {
		return "", err
	}
	return id, nil
}

func (u *Usecase) ListSubjects(ctx context.Context, page, pageSize int64) ([]q.Subject, error) {
	tid, err := u.tenantFromContext(ctx)
	if err != nil {
		return nil, err
	}
	limit, offset := normalize(page, pageSize)
	return u.Subjects.ListByTenant(ctx, tid, limit, offset)
}

func (u *Usecase) UpdateSubject(ctx context.Context, req schooldto.UpdateSubjectRequest) error {
	cur, err := u.Subjects.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}
	newName := cur.Name
	if req.Name != nil {
		newName = *req.Name
	}
	var dep sql.NullString
	if req.DepartmentID != nil {
		dep = sql.NullString{String: *req.DepartmentID, Valid: true}
	}
	return u.Subjects.Update(ctx, q.UpdateSubjectParams{
		ID:           req.ID,
		Name:         newName,
		DepartmentID: dep,
	})
}

func (u *Usecase) DeleteSubject(ctx context.Context, id string) error {
	return u.Subjects.Delete(ctx, id)
}

// Teachers
func (u *Usecase) CreateTeacher(ctx context.Context, req schooldto.CreateTeacherRequest) (*schooldto.CreateTeacherResponse, error) {
	tid, err := u.tenantFromContext(ctx)
	if err != nil {
		return nil, err
	}
	base := strings.ToLower(string([]rune(req.FirstName)[0]) + req.LastName)
	username := fmt.Sprintf("%s%s", base, randHex(4))
	for {
		exists, err := u.Users.GetByUsername(ctx, username)
		if err != nil {
			return nil, err
		}
		if exists == nil {
			break
		}
		username = fmt.Sprintf("%s%s", base, randHex(4))
	}
	password := randHex(6)
	hash, err := security.HashPassword(password)
	if err != nil {
		return nil, err
	}
	roleRow, err := u.Queries.GetRoleByName(ctx, "teacher")
	if err != nil {
		return nil, err
	}
	roleID := roleRow.ID
	userID := uuid.New().String()
	if err := u.Users.Create(ctx, domain.User{
		ID:           userID,
		TenantID:     &tid,
		Username:     username,
		PasswordHash: hash,
		Email:        req.Email,
		RoleID:       &roleID,
		Status:       "active",
	}); err != nil {
		return nil, err
	}
	tID := uuid.New().String()
	var dep sql.NullString
	if req.DepartmentID != nil {
		dep = sql.NullString{String: *req.DepartmentID, Valid: true}
	}
	if err := u.Teachers.Create(ctx, q.CreateTeacherParams{
		ID:           tID,
		TenantID:     sql.NullString{String: tid, Valid: true},
		FirstName:    sql.NullString{String: req.FirstName, Valid: true},
		LastName:     sql.NullString{String: req.LastName, Valid: true},
		DepartmentID: dep,
		UserID:       sql.NullString{String: userID, Valid: true},
	}); err != nil {
		return nil, err
	}
	return &schooldto.CreateTeacherResponse{
		TeacherID: tID,
		Username:  username,
		Password:  password,
	}, nil
}

func (u *Usecase) ListTeachers(ctx context.Context, page, pageSize int64) ([]q.Teacher, error) {
	tid, err := u.tenantFromContext(ctx)
	if err != nil {
		return nil, err
	}
	limit, offset := normalize(page, pageSize)
	return u.Teachers.ListByTenant(ctx, tid, limit, offset)
}

func (u *Usecase) UpdateTeacher(ctx context.Context, req schooldto.UpdateTeacherRequest) error {
	return u.Teachers.Update(ctx, q.UpdateTeacherParams{
		ID:           req.ID,
		FirstName:    toNull(req.FirstName),
		LastName:     toNull(req.LastName),
		DepartmentID: toNull(req.DepartmentID),
	})
}

func (u *Usecase) DeleteTeacher(ctx context.Context, id string) error {
	return u.Teachers.Delete(ctx, id)
}

func (u *Usecase) Assign(ctx context.Context, req schooldto.AssignRequest) error {
	return u.Teachers.Assign(ctx, req.TeacherID, req.SubjectID, req.SectionID)
}

func (u *Usecase) Unassign(ctx context.Context, req schooldto.UnassignRequest) error {
	return u.Teachers.Unassign(ctx, req.TeacherID, req.SubjectID, req.SectionID)
}

func (u *Usecase) ListTeacherAssignments(ctx context.Context, page, pageSize int64) (interface{}, error) {
	tenantID, err := u.tenantFromContext(ctx)
	if err != nil {
		return nil, err
	}
	limit, offset := normalize(page, pageSize)

	return u.Teachers.ListAssignmentsDetailedByTenant(ctx, tenantID, limit, offset)
}

func randHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func normalize(page, pageSize int64) (int64, int64) {
	if pageSize <= 0 {
		pageSize = 20
	}
	if page <= 0 {
		page = 1
	}
	return pageSize, (page - 1) * pageSize
}

func toNull(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{}
}
