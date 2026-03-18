package superadmin

import (
	"context"
	"crypto/rand"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"database/sql"
	"school-exam/internal/domain"
	sdto "school-exam/internal/dto/superadmin"
	"school-exam/internal/repository"
	"school-exam/internal/security"
	q "school-exam/internal/sqlc/gen"

	"github.com/google/uuid"
)

type Usecase struct {
	Tenants  *repository.TenantRepository
	Users    *repository.UserRepository
	Students *repository.StudentRepository
	Sections *repository.SectionRepository
	Queries  *q.Queries
}

func NewUsecase(t *repository.TenantRepository, u *repository.UserRepository, s *repository.StudentRepository, sec *repository.SectionRepository) *Usecase {
	return &Usecase{Tenants: t, Users: u, Students: s, Sections: sec, Queries: q.New(t.DB)}
}

func (u *Usecase) CreateTenant(ctx context.Context, req sdto.CreateTenantRequest) (*sdto.CreateTenantResponse, error) {
	id := uuid.New().String()
	if err := u.Tenants.Create(ctx, id, req.Name, req.Address, req.Phone); err != nil {
		return nil, err
	}
	// create a school_admin for this tenant
	// username: 'admin' + 4 hex
	base := "admin"
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
	roleRow, err := u.Queries.GetRoleByName(ctx, "school_admin")
	if err != nil {
		return nil, err
	}
	roleID := roleRow.ID
	if err := u.Users.Create(ctx, domain.User{
		ID:           uuid.New().String(),
		TenantID:     &id,
		Username:     username,
		PasswordHash: hash,
		RoleID:       &roleID,
		Status:       "active",
	}); err != nil {
		return nil, err
	}
	return &sdto.CreateTenantResponse{ID: id, AdminUsername: username, AdminPassword: password}, nil
}

func (u *Usecase) UpdateTenant(ctx context.Context, req sdto.UpdateTenantRequest) error {
	return u.Tenants.Update(ctx, req.ID, req.Name, req.Address, req.Phone)
}

func (u *Usecase) ActivateTenant(ctx context.Context, id string) error {
	return u.Tenants.SetStatus(ctx, id, "active")
}

func (u *Usecase) DeactivateTenant(ctx context.Context, id string) error {
	return u.Tenants.SetStatus(ctx, id, "inactive")
}

func (u *Usecase) GetTenant(ctx context.Context, id string) (q.Tenant, error) {
	return u.Tenants.GetByID(ctx, id)
}

func (u *Usecase) ListTenants(ctx context.Context, status *string, page, pageSize int64) ([]q.Tenant, error) {
	limit, offset := normalizePagination(page, pageSize)
	if status != nil && *status != "" {
		return u.Tenants.ListByStatus(ctx, *status, limit, offset)
	}
	return u.Tenants.List(ctx, limit, offset)
}

func (u *Usecase) ListSections(ctx context.Context, tenantID string, page, pageSize int64) ([]q.Section, error) {
	limit, offset := normalizePagination(page, pageSize)
	return u.Sections.ListByTenant(ctx, tenantID, limit, offset)
}

func (u *Usecase) DeleteTenant(ctx context.Context, id string) error {
	return u.Tenants.SetStatus(ctx, id, "inactive")
}

func (u *Usecase) UpdateStudent(ctx context.Context, req sdto.UpdateStudentRequest) error {
	if req.Status != nil {
		if err := u.Students.SetStatus(ctx, req.ID, *req.Status); err != nil {
			return err
		}
	}
	return u.Students.Update(ctx, req.ID, req.FirstName, req.LastName, req.Year, req.SectionID, req.DepartmentID)
}

func (u *Usecase) CreateStudent(ctx context.Context, req sdto.CreateStudentRequest) (*sdto.CreateStudentResponse, error) {
	base := strings.ToLower(string([]rune(req.FirstName)[0]) + req.LastName)
	suffix := randHex(4)
	username := fmt.Sprintf("%s%s", base, suffix)
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
	roleRow, err := u.Queries.GetRoleByName(ctx, "student")
	if err != nil {
		return nil, err
	}
	roleID := roleRow.ID
	userID := uuid.New().String()
	dUser := domain.User{
		ID:           userID,
		TenantID:     &req.TenantID,
		Username:     username,
		PasswordHash: hash,
		Email:        req.Email,
		RoleID:       &roleID,
		Status:       "active",
	}
	if err := u.Users.Create(ctx, dUser); err != nil {
		return nil, err
	}
	stID := uuid.New().String()
	var section, dept, year sql.NullString
	if req.SectionID != nil {
		section = sql.NullString{Valid: true, String: *req.SectionID}
	}
	if req.DepartmentID != nil {
		dept = sql.NullString{Valid: true, String: *req.DepartmentID}
	}
	if req.Year != nil {
		year = sql.NullString{Valid: true, String: *req.Year}
	}
	err = u.Students.Create(ctx, q.CreateStudentParams{
		ID:           stID,
		TenantID:     sql.NullString{Valid: true, String: req.TenantID},
		StudentCode:  sql.NullString{Valid: true, String: req.StudentCode},
		FirstName:    sql.NullString{Valid: true, String: req.FirstName},
		LastName:     sql.NullString{Valid: true, String: req.LastName},
		Year:         year,
		SectionID:    section,
		DepartmentID: dept,
		UserID:       sql.NullString{Valid: true, String: userID},
	})
	if err != nil {
		return nil, err
	}
	return &sdto.CreateStudentResponse{
		StudentID: stID,
		Username:  username,
		Password:  password,
	}, nil
}

type ImportResult struct {
	Row         int    `json:"row"`
	StudentCode string `json:"student_code"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Error       string `json:"error"`
}

func (u *Usecase) ImportStudents(ctx context.Context, tenantID string, r io.Reader) ([]ImportResult, error) {
	cr := csv.NewReader(r)
	header, err := cr.Read()
	if err != nil {
		return nil, err
	}
	col := map[string]int{}
	for i, h := range header {
		col[strings.ToLower(strings.TrimSpace(h))] = i
	}
	var res []ImportResult
	rowNum := 1
	for {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}
		rowNum++
		if err != nil {
			res = append(res, ImportResult{Row: rowNum, Error: err.Error()})
			continue
		}
		get := func(name string) string {
			if idx, ok := col[name]; ok && idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
			return ""
		}
		req := sdto.CreateStudentRequest{
			TenantID:    tenantID,
			StudentCode: get("student_code"),
			FirstName:   get("first_name"),
			LastName:    get("last_name"),
		}
		sec := get("section_id")
		if sec != "" {
			req.SectionID = &sec
		}
		dep := get("department_id")
		if dep != "" {
			req.DepartmentID = &dep
		}
		em := get("email")
		if em != "" {
			req.Email = &em
		}
		yr := get("year")
		if yr != "" {
			req.Year = &yr
		}
		out, err := u.CreateStudent(ctx, req)
		if err != nil {
			res = append(res, ImportResult{
				Row:         rowNum,
				StudentCode: req.StudentCode,
				Error:       err.Error(),
			})
			continue
		}
		res = append(res, ImportResult{
			Row:         rowNum,
			StudentCode: req.StudentCode,
			Username:    out.Username,
			Password:    out.Password,
		})
	}
	return res, nil
}

func randHex(nBytes int) string {
	b := make([]byte, nBytes)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func normalizePagination(page, pageSize int64) (int64, int64) {
	if pageSize <= 0 {
		pageSize = 20
	}
	if page <= 0 {
		page = 1
	}
	limit := pageSize
	offset := (page - 1) * pageSize
	return limit, offset
}
