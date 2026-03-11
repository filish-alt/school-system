package route

import (
	"time"

	hAuth "school-exam/internal/handler/auth"
	hSchool "school-exam/internal/handler/school"
	hSuper "school-exam/internal/handler/superadmin"
	hTeacher "school-exam/internal/handler/teacher"
	hExam "school-exam/internal/handler/exam"
	hSession "school-exam/internal/handler/exam_session"
	"school-exam/internal/middleware"
	"school-exam/internal/module/auth"
	"school-exam/internal/module/exam"
	"school-exam/internal/module/exam_session"
	"school-exam/internal/module/school"
	"school-exam/internal/module/superadmin"
	"school-exam/internal/module/teacher"
	"school-exam/internal/security"

	"github.com/gin-gonic/gin"
)

func SetupRouter(authUC *auth.Usecase, superUC *superadmin.Usecase, schoolUC *school.Usecase, teacherUC *teacher.Usecase, examUC *exam.Usecase, sessionUC *exam_session.Usecase, ts security.TokenService) *gin.Engine {
	r := gin.Default()
	r.Use(CORSMiddleware())
	
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	ah := hAuth.New(authUC)
	// public group
	v1 := r.Group("/api/v1")
	v1.POST("/auth/login", ah.Login)
	
	sessH := hSession.NewHandler(sessionUC)
	// authenticated group (same prefix, adds auth)
	authed := r.Group("/api/v1")
	authed.Use(middleware.Auth(ts))
	authed.PATCH("/auth/password", ah.UpdatePassword)

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
	schoolGroup.GET("/assignments", sch.ListAssignedTeachers)
	schoolGroup.GET("/violations", sessH.ListAllViolations)
	
	// teacher routes
	teacherGroup := authed.Group("/teacher")
	teacherGroup.Use(middleware.RequireRoles("teacher"))
	th := hTeacher.New(teacherUC)
	teacherGroup.POST("/question-banks", th.CreateQuestionBank)
	teacherGroup.GET("/question-banks", th.ListQuestionBanks)
	teacherGroup.POST("/questions", th.CreateQuestion)
	teacherGroup.GET("/questions", th.ListQuestions)
	teacherGroup.POST("/questions/import", th.ImportQuestions)
	teacherGroup.PATCH("/questions", th.UpdateQuestion)
	teacherGroup.DELETE("/questions/:id", th.DeleteQuestion)
	teacherGroup.POST("/options", th.CreateOption)
	teacherGroup.GET("/options", th.ListOptions)
	teacherGroup.PATCH("/options", th.UpdateOption)
	teacherGroup.DELETE("/options/:id", th.DeleteOption)
	teacherGroup.GET("/students", th.ListMyStudents)
	teacherGroup.GET("/my-assignments", th.ListMyAssignments)

	// exam routes
	eh := hExam.New(examUC)
	teacherGroup.POST("/exams", eh.CreateExam)
	teacherGroup.GET("/exams", eh.ListExams)
	teacherGroup.GET("/exams/:id", eh.GetExam)
	teacherGroup.GET("/exams/:id/questions", eh.GetTeacherExamQuestions)
	teacherGroup.PATCH("/exams", eh.UpdateExam)
	teacherGroup.PATCH("/exams/:id/status", eh.UpdateStatus)
	teacherGroup.DELETE("/exams/:id", eh.DeleteExam)
	teacherGroup.POST("/exams/questions", eh.AddQuestions)
	teacherGroup.POST("/exams/questions/random", eh.AddRandomQuestions)
	teacherGroup.DELETE("/exams/questions/:id", eh.RemoveQuestion)
	teacherGroup.GET("/exams/:id/marks", eh.GetExamMarks)
	teacherGroup.GET("/exams/:id/marks/download", eh.DownloadExamMarks)

	// student routes
	studentGroup := authed.Group("/student")
	studentGroup.Use(middleware.RequireRoles("student"))
	{
		studentGroup.GET("/exams", eh.ListStudentExams)
		studentGroup.GET("/exams/:id", eh.GetStudentExam)
		studentGroup.GET("/sessions", sessH.ListMySessions)
		studentGroup.POST("/sessions/violations", sessH.ReportViolation)
		studentGroup.POST("/sessions/start", sessH.StartSession)
		studentGroup.POST("/sessions/answers", sessH.SaveAnswer)
		studentGroup.POST("/sessions/submit", sessH.SubmitSession)
		studentGroup.GET("/sessions/:id", sessH.GetSession)
		studentGroup.GET("/sessions/:id/result", sessH.GetSessionResult)
	}

	authed.GET("/me", func(c *gin.Context) { c.JSON(200, gin.H{"time": time.Now().UTC()}) })

	return r
}


func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control,X-API-Key, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}