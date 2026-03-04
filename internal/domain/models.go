package domain

type Role struct {
	ID   int64
	Name string
}

type User struct {
	ID           string
	TenantID     *string
	Username     string
	PasswordHash string
	Email        *string
	RoleID       *int64
	RoleName     *string
	Status       string
}

