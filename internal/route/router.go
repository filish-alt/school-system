package route

import (
	"time"

	hAuth "school-exam/internal/handler/auth"
	hSchool "school-exam/internal/handler/school"
	hSuper "school-exam/internal/handler/superadmin"
	"school-exam/internal/middleware"
	"school-exam/internal/module/auth"
	"school-exam/internal/module/school"
	"school-exam/internal/module/superadmin"
	"school-exam/internal/security"

	"github.com/gin-gonic/gin"
)

func SetupRouter(authUC *auth.Usecase, superUC *superadmin.Usecase, schoolUC *school.Usecase, ts security.TokenService) *gin.Engine {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	ah := hAuth.New(authUC)
	// public group
	v1 := r.Group("/api/v1")
	v1.POST("/auth/login", ah.Login)
	// authenticated group (same prefix, adds auth)
	authed := r.Group("/api/v1")
	authed.Use(middleware.Auth(ts))

	// super admin routes (must be authenticated)
	sAdmin := authed.Group("/superadmin")
	sAdmin.Use(middleware.RequireRoles("super_admin"))

	sh := hSuper.New(superUC)
	sAdmin.GET("/tenants", sh.ListTenants)
	sAdmin.GET("/tenants/:id", sh.GetTenant)
	sAdmin.POST("/tenants", sh.CreateTenant)
	sAdmin.PATCH("/tenants", sh.UpdateTenant)
	sAdmin.DELETE("/tenants/:id", sh.DeleteTenant)
	sAdmin.POST("/tenants/:id/activate", sh.ActivateTenant)
	sAdmin.POST("/tenants/:id/deactivate", sh.DeactivateTenant)
	sAdmin.GET("/students", sh.ListStudents)
	sAdmin.GET("/students/:id", sh.GetStudent)
	sAdmin.POST("/students", sh.CreateStudent)
	sAdmin.PATCH("/students/:id", sh.UpdateStudent)
	sAdmin.DELETE("/students/:id", sh.DeleteStudent)
	sAdmin.POST("/students/import", sh.ImportStudents)

	// school admin routes (must be authenticated)
	schoolGroup := authed.Group("/school")
	schoolGroup.Use(middleware.RequireRoles("school_admin"))
	sch := hSchool.New(schoolUC)
	schoolGroup.POST("/departments", sch.CreateDepartment)
	schoolGroup.GET("/departments", sch.ListDepartments)
	schoolGroup.PATCH("/departments", sch.UpdateDepartment)
	schoolGroup.DELETE("/departments/:id", sch.DeleteDepartment)
	schoolGroup.POST("/sections", sch.CreateSection)
	schoolGroup.GET("/sections", sch.ListSections)
	schoolGroup.PATCH("/sections", sch.UpdateSection)
	schoolGroup.DELETE("/sections/:id", sch.DeleteSection)
	schoolGroup.POST("/subjects", sch.CreateSubject)
	schoolGroup.GET("/subjects", sch.ListSubjects)
	schoolGroup.PATCH("/subjects", sch.UpdateSubject)
	schoolGroup.DELETE("/subjects/:id", sch.DeleteSubject)
	schoolGroup.POST("/teachers", sch.CreateTeacher)
	schoolGroup.GET("/teachers", sch.ListTeachers)
	schoolGroup.PATCH("/teachers", sch.UpdateTeacher)
	schoolGroup.DELETE("/teachers/:id", sch.DeleteTeacher)
	schoolGroup.POST("/assignments", sch.Assign)
	schoolGroup.DELETE("/assignments", sch.Unassign)
	authed.GET("/me", func(c *gin.Context) { c.JSON(200, gin.H{"time": time.Now().UTC()}) })
	return r
}
