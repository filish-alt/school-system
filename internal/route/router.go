package route

import (
	"time"

	"school-exam/internal/config"
	hAcademic "school-exam/internal/handler/academic"
	hAuth "school-exam/internal/handler/auth"
	hComplaint "school-exam/internal/handler/complaint"
	hCoursework "school-exam/internal/handler/coursework"
	hExam "school-exam/internal/handler/exam"
	hExamControl "school-exam/internal/handler/exam_control"
	hSession "school-exam/internal/handler/exam_session"
	hResult "school-exam/internal/handler/result"
	hSchedule "school-exam/internal/handler/schedule"
	hSchool "school-exam/internal/handler/school"
	hSuper "school-exam/internal/handler/superadmin"
	hSystem "school-exam/internal/handler/system"
	hTeacher "school-exam/internal/handler/teacher"
	"school-exam/internal/middleware"
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
	"school-exam/internal/security"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authUC *auth.Usecase,
	superUC *superadmin.Usecase,
	schoolUC *school.Usecase,
	teacherUC *teacher.Usecase,
	examUC *exam.Usecase,
	sessionUC *exam_session.Usecase,
	acadUC *academic.Usecase,
	cwUC *coursework.Usecase,
	ecUC *exam_control.Usecase,
	compUC *complaint.Usecase,
	schUC *schedule.Usecase,
	resUC *result.Usecase,
	sysUC *system.Usecase,
	ts security.TokenService,
	cfg config.Config,
) *gin.Engine {
	r := gin.Default()
	r.Use(CORSMiddleware())

	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	r.Static("/uploads", "./uploads")

	ah := hAuth.New(authUC)
	sessH := hSession.NewHandler(sessionUC)
	sh := hSuper.New(superUC)
	sch := hSchool.New(schoolUC)
	th := hTeacher.New(teacherUC)
	eh := hExam.New(examUC)

	acadH := hAcademic.New(acadUC)
	cwH := hCoursework.New(cwUC)
	ecH := hExamControl.New(ecUC)
	compH := hComplaint.New(compUC)
	schH := hSchedule.New(schUC)
	resH := hResult.New(resUC)
	sysH := hSystem.New(sysUC, cfg)

	// Public routes
	v1 := r.Group("/api/v1")
	v1.POST("/auth/login", ah.Login)
	v1.GET("/system/emergency-exam/:examId/print", sysH.PrintEmergencyExamHTML)
	v1.GET("/settings", acadH.GetSettings)

	// Authenticated routes
	authed := r.Group("/api/v1")
	authed.Use(middleware.Auth(ts))
	authed.PATCH("/auth/password", ah.UpdatePassword)
	authed.GET("/me", func(c *gin.Context) { c.JSON(200, gin.H{"time": time.Now().UTC()}) })
	authed.POST("/settings", acadH.SaveSetting)

	// 1. SUPER ADMIN
	sAdmin := authed.Group("/superadmin")
	sAdmin.Use(middleware.RequireRoles("super_admin"))
	{
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
		sAdmin.GET("/sections", sh.ListSections)
		sAdmin.POST("/upload", th.UploadImage)

		// Academic Years & Settings
		sAdmin.POST("/academic-years", acadH.CreateAcademicYear)
		sAdmin.GET("/academic-years", acadH.ListAcademicYears)
		sAdmin.POST("/academic-years/:id/activate", acadH.SetActiveAcademicYear)
		sAdmin.DELETE("/academic-years/:id", acadH.DeleteAcademicYear)
		sAdmin.GET("/settings", acadH.GetSettings)
		sAdmin.POST("/settings", acadH.SaveSetting)

		// Audit Logs & Backups
		sAdmin.GET("/audit-logs", sysH.ListAuditLogs)
		sAdmin.POST("/backups", sysH.CreateBackup)
		sAdmin.GET("/backups", sysH.ListBackups)
		sAdmin.GET("/backups/:filename/download", sysH.DownloadBackup)
	}

	// 2. SCHOOL DIRECTORS, ADMINS, DEPT HEADS & TEACHERS
	schoolGroup := authed.Group("/school")
	schoolGroup.Use(middleware.RequireRoles("super_admin", "school_admin", "school_director", "vice_director", "department_head", "teacher"))
	{
		// Operational Academic Structure (Creation & Management restricted to School Leadership & Admins)
		acadStructureAdmin := schoolGroup.Group("")
		acadStructureAdmin.Use(middleware.RequireRoles("super_admin", "school_admin", "school_director", "vice_director"))
		{
			acadStructureAdmin.POST("/departments", sch.CreateDepartment)
			acadStructureAdmin.PATCH("/departments", sch.UpdateDepartment)
			acadStructureAdmin.DELETE("/departments/:id", sch.DeleteDepartment)

			acadStructureAdmin.POST("/sections", sch.CreateSection)
			acadStructureAdmin.PATCH("/sections", sch.UpdateSection)
			acadStructureAdmin.DELETE("/sections/:id", sch.DeleteSection)

			acadStructureAdmin.POST("/subjects", sch.CreateSubject)
			acadStructureAdmin.PATCH("/subjects", sch.UpdateSubject)
			acadStructureAdmin.DELETE("/subjects/:id", sch.DeleteSubject)
		}

		// Read-only listing for all authorized academic staff
		schoolGroup.GET("/departments", sch.ListDepartments)
		schoolGroup.GET("/sections", sch.ListSections)
		schoolGroup.GET("/subjects", sch.ListSubjects)

		// Staff / Faculty creation & deletion (Restricted to Super Admin & School Leadership)
		staffAdmin := schoolGroup.Group("")
		staffAdmin.Use(middleware.RequireRoles("super_admin", "school_admin", "school_director", "vice_director"))
		{
			staffAdmin.POST("/teachers", sch.CreateTeacher)
			staffAdmin.PATCH("/teachers", sch.UpdateTeacher)
			staffAdmin.DELETE("/teachers/:id", sch.DeleteTeacher)
		}
		schoolGroup.GET("/teachers", sch.ListTeachers)

		// Teacher Section & Class Assignments (Department Head & School Leadership Task)
		assignAdmin := schoolGroup.Group("")
		assignAdmin.Use(middleware.RequireRoles("super_admin", "school_admin", "school_director", "vice_director", "department_head"))
		{
			assignAdmin.POST("/assignments", sch.Assign)
			assignAdmin.DELETE("/assignments", sch.Unassign)
			assignAdmin.GET("/assignments", sch.ListAssignedTeachers)
		}
		schoolGroup.GET("/violations", sessH.ListAllViolations)

		// Timetable Schedules
		schoolGroup.POST("/schedules", schH.CreateSchedule)
		schoolGroup.GET("/schedules/teacher/:teacherId", schH.ListByTeacher)
		schoolGroup.GET("/schedules/section/:sectionId", schH.ListBySection)
		schoolGroup.DELETE("/schedules/:id", schH.DeleteSchedule)

		// Results & Transfer Sheets
		schoolGroup.POST("/results/calculate", resH.CalculateResults)
		schoolGroup.GET("/results/transfer-sheet/:sectionId", resH.GetTransferSheet)

		// Complaints resolution
		schoolGroup.GET("/complaints", compH.ListSchoolComplaints)
		schoolGroup.PATCH("/complaints/:id/resolve", compH.ResolveComplaint)

		// Academic years list
		schoolGroup.GET("/academic-years", acadH.ListAcademicYears)

		// Tenant Assessment & Grading Settings
		schoolGroup.GET("/settings", acadH.GetSettings)
		schoolGroup.POST("/settings", acadH.SaveSetting)
		schoolGroup.POST("/upload", th.UploadImage)
	}

	// 3. TEACHER
	teacherGroup := authed.Group("/teacher")
	teacherGroup.Use(middleware.RequireRoles("teacher", "department_head"))
	{
		// Question Banks & Questions
		teacherGroup.POST("/question-banks", th.CreateQuestionBank)
		teacherGroup.GET("/question-banks", th.ListQuestionBanks)
		teacherGroup.POST("/upload", th.UploadImage)
		teacherGroup.POST("/questions", th.CreateQuestion)
		teacherGroup.GET("/questions", th.ListQuestions)
		teacherGroup.POST("/questions/import", th.ImportQuestions)
		teacherGroup.PATCH("/questions", th.UpdateQuestion)
		teacherGroup.DELETE("/questions/:id", th.DeleteQuestion)
		teacherGroup.POST("/options", th.CreateOption)
		teacherGroup.GET("/options", th.ListOptions)
		teacherGroup.PATCH("/options", th.UpdateOption)
		teacherGroup.DELETE("/options/:id", th.DeleteOption)
		teacherGroup.POST("/questions/:questionId/match-pairs", ecH.SetMatchPairs)
		teacherGroup.GET("/questions/:questionId/match-pairs", ecH.GetMatchPairs)

		teacherGroup.GET("/students", th.ListMyStudents)
		teacherGroup.GET("/my-assignments", th.ListMyAssignments)
		teacherGroup.GET("/my-schedule", schH.ListMySchedule)

		// Continuous Assessment / Coursework
		teacherGroup.POST("/coursework", cwH.CreateCoursework)
		teacherGroup.GET("/coursework", cwH.ListBySectionSubject)
		teacherGroup.GET("/my-coursework", cwH.ListMyCourseworks)
		teacherGroup.POST("/coursework/marks", cwH.SaveMark)
		teacherGroup.GET("/coursework/:courseworkId/marks", cwH.ListMarksByCoursework)

		// Exams & Question linking
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

		// Live Exam Controls
		teacherGroup.PATCH("/exams/:id/controls", ecH.UpdateControls)
		teacherGroup.POST("/exams/extensions", ecH.AddExtension)
		teacherGroup.GET("/exams/:id/extensions", ecH.ListExtensions)
		teacherGroup.POST("/exams/remote-authorizations", ecH.AddRemoteAuth)
		teacherGroup.GET("/exams/:id/remote-authorizations", ecH.ListRemoteAuths)
		teacherGroup.POST("/exams/drop-question", ecH.DropQuestion)
		teacherGroup.POST("/exams/correct-key", ecH.CorrectAnswerKey)

		// Emergency Exam Package & Printable Fallback
		teacherGroup.GET("/exams/:id/emergency-package", sysH.GetEmergencyExamPackage)

		// Results calculation trigger
		teacherGroup.POST("/results/calculate", resH.CalculateResults)
		teacherGroup.GET("/results/transfer-sheet/:sectionId", resH.GetTransferSheet)
	}

	// 4. STUDENT
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

		// Continuous Assessment & Final Results
		studentGroup.GET("/my-marks", cwH.ListMyMarks)
		studentGroup.GET("/my-results", resH.GetMyResults)

		// Complaints
		studentGroup.POST("/complaints", compH.CreateComplaint)
		studentGroup.GET("/my-complaints", compH.ListMyComplaints)
	}

	return r
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-API-Key, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}