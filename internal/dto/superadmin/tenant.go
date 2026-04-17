package superadmindto

type CreateTenantRequest struct {
	Name          string  `json:"name" binding:"required"`
	Address       *string `json:"address"`
	Phone         *string `json:"phone"`
	AdminUsername string  `json:"admin_username" binding:"required"`
}

type CreateTenantResponse struct {
	ID            string `json:"id"`
	AdminUsername string `json:"admin_username"`
	AdminPassword string `json:"admin_password"`
}

type UpdateTenantRequest struct {
	ID      string  `json:"id" binding:"required"`
	Name    string  `json:"name" binding:"required"`
	Address *string `json:"address"`
	Phone   *string `json:"phone"`
}

type TenantListQuery struct {
	Page     int64 `form:"page"`
	PageSize int64 `form:"page_size"`
}
