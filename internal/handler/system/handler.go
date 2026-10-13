package system

import (
	"fmt"
	"net/http"
	"path/filepath"
	"school-exam/internal/config"
	"school-exam/internal/middleware"
	"school-exam/internal/module/system"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc  *system.Usecase
	cfg config.Config
}

func New(uc *system.Usecase, cfg config.Config) *Handler {
	return &Handler{uc: uc, cfg: cfg}
}

func (h *Handler) ListAuditLogs(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	list, err := h.uc.ListAuditLogs(c.Request.Context(), tenantID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) CreateBackup(c *gin.Context) {
	backupPath, err := h.uc.CreateBackup(c.Request.Context(), h.cfg.DBPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "backup created", "path": backupPath})
}

func (h *Handler) ListBackups(c *gin.Context) {
	list, err := h.uc.ListBackups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) DownloadBackup(c *gin.Context) {
	filename := filepath.Base(c.Param("filename"))
	filePath := filepath.Join("backups", filename)
	c.FileAttachment(filePath, filename)
}

func (h *Handler) GetEmergencyExamPackage(c *gin.Context) {
	examID := c.Param("id")
	if examID == "" {
		examID = c.Param("examId")
	}
	pkg, err := h.uc.GenerateEmergencyPackage(c.Request.Context(), examID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pkg)
}

func (h *Handler) PrintEmergencyExamHTML(c *gin.Context) {
	examID := c.Param("id")
	if examID == "" {
		examID = c.Param("examId")
	}
	pkg, err := h.uc.GenerateEmergencyPackage(c.Request.Context(), examID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<title>%s - Exam Paper & Emergency Package</title>
	<style>
		body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; margin: 40px; color: #111; line-height: 1.5; }
		.header { text-align: center; border-bottom: 2px solid #333; padding-bottom: 12px; margin-bottom: 20px; }
		.meta-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; margin-bottom: 20px; font-size: 14px; background: #f8f9fa; padding: 12px; border-radius: 4px; }
		.instructions { background: #fff3cd; border-left: 4px solid #ffeeba; padding: 10px; margin-bottom: 25px; font-size: 13px; }
		.question { margin-bottom: 24px; page-break-inside: avoid; }
		.question-title { font-weight: 600; margin-bottom: 8px; }
		.options { margin-left: 20px; list-style-type: upper-alpha; }
		.options li { margin-bottom: 4px; }
		.page-break { page-break-before: always; }
		table { width: 100%%; border-collapse: collapse; margin-top: 15px; font-size: 13px; }
		th, td { border: 1px solid #ccc; padding: 8px 12px; text-align: left; }
		th { background: #f0f0f0; }
		@media print {
			body { margin: 20mm 15mm; }
			.no-print { display: none; }
		}
	</style>
</head>
<body>
	<div class="no-print" style="margin-bottom: 20px; text-align: right;">
		<button onclick="window.print()" style="padding: 8px 16px; background: #2563eb; color: #fff; border: none; border-radius: 4px; cursor: pointer; font-weight: bold;">🖨️ Print Exam Package</button>
	</div>

	<!-- SECTION 1: QUESTION PAPER -->
	<div class="header">
		<h2 style="margin: 0;">ACADEMIC EXAMINATION PAPER</h2>
		<h3 style="margin: 4px 0 0 0; color: #4b5563;">%s</h3>
	</div>

	<div class="meta-grid">
		<div><strong>Subject:</strong> %s</div>
		<div><strong>Grade / Section:</strong> %s</div>
		<div><strong>Duration:</strong> %d Minutes</div>
		<div><strong>Total Marks:</strong> %.1f</div>
		<div><strong>Academic Year:</strong> %s</div>
		<div><strong>Date:</strong> %s</div>
	</div>

	<div class="instructions">
		<strong>Instructions to Students:</strong> %s
	</div>

	<div style="margin-bottom: 15px; font-weight: bold; border-bottom: 1px solid #ddd; padding-bottom: 4px;">QUESTIONS</div>
`, pkg.ExamTitle, pkg.ExamTitle, pkg.SubjectName, pkg.SectionName, pkg.DurationMinutes, pkg.TotalMarks, pkg.AcademicYear, pkg.GeneratedAt.Format("2006-01-02"), pkg.Instructions)

	for _, q := range pkg.Questions {
		html += fmt.Sprintf(`
	<div class="question">
		<div class="question-title">%d. %s <span style="float: right; color: #666; font-size: 12px;">[%.1f Marks]</span></div>`, q.QuestionNumber, q.QuestionText, q.Marks)
		if len(q.Options) > 0 {
			html += `<ol class="options">`
			for _, opt := range q.Options {
				html += fmt.Sprintf(`<li>%s</li>`, opt)
			}
			html += `</ol>`
		} else {
			html += `<div style="height: 60px; border-bottom: 1px dashed #ccc; margin-top: 10px;"></div>`
		}
		html += `</div>`
	}

	// SECTION 2: ATTENDANCE & MARK-ENTRY SHEET
	html += fmt.Sprintf(`
	<div class="page-break"></div>
	<div class="header">
		<h2 style="margin: 0;">EXAMINATION ATTENDANCE & MARK-ENTRY SHEET</h2>
		<h4 style="margin: 4px 0 0 0; color: #4b5563;">%s - %s (%s)</h4>
	</div>
	<table>
		<thead>
			<tr>
				<th style="width: 40px;">#</th>
				<th style="width: 120px;">Student Code</th>
				<th>Full Name</th>
				<th style="width: 140px;">Student Signature</th>
				<th style="width: 100px;">Marks Obtained</th>
				<th>Invigilator Remarks</th>
			</tr>
		</thead>
		<tbody>`, pkg.ExamTitle, pkg.SubjectName, pkg.SectionName)

	for idx, s := range pkg.StudentRoster {
		html += fmt.Sprintf(`
			<tr>
				<td>%d</td>
				<td>%s</td>
				<td>%s</td>
				<td></td>
				<td></td>
				<td></td>
			</tr>`, idx+1, s.StudentCode, s.FullName)
	}

	html += `
		</tbody>
	</table>

	<div style="margin-top: 40px; display: flex; justify-content: space-between; font-size: 13px;">
		<div>Invigilator Name & Signature: _______________________</div>
		<div>Date: _______________</div>
	</div>
</body>
</html>`

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}
