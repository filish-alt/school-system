package exam

import (
	"net/http"

	"github.com/gin-gonic/gin"
	examdto "school-exam/internal/dto/exam"
	"school-exam/internal/module/exam"
)

type Handler struct {
	UC *exam.Usecase
}

func New(uc *exam.Usecase) *Handler {
	return &Handler{UC: uc}
}

func (h *Handler) CreateExam(c *gin.Context) {
	var req examdto.CreateExamRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	id, err := h.UC.CreateExam(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) GetExam(c *gin.Context) {
	id := c.Param("id")
	out, err := h.UC.GetExam(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) GetTeacherExamQuestions(c *gin.Context) {
	id := c.Param("id")
	out, err := h.UC.GetTeacherExamQuestions(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) ListExams(c *gin.Context) {
	var q examdto.ListQuery
	_ = c.ShouldBindQuery(&q)
	out, err := h.UC.ListExams(c, q.Page, q.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) ListStudentExams(c *gin.Context) {
	out, err := h.UC.ListStudentExams(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) GetStudentExam(c *gin.Context) {
	id := c.Param("id")
	out, err := h.UC.GetStudentExam(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) UpdateExam(c *gin.Context) {
	var req examdto.UpdateExamRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.UC.UpdateExam(c, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var req examdto.UpdateStatusRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.UC.UpdateStatus(c, id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) DeleteExam(c *gin.Context) {
	id := c.Param("id")
	if err := h.UC.DeleteExam(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) AddQuestions(c *gin.Context) {
	var req examdto.AddQuestionsRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.UC.AddQuestions(c, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "questions added"})
}

func (h *Handler) AddRandomQuestions(c *gin.Context) {
	var req examdto.AddRandomQuestionsRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.UC.AddRandomQuestions(c, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "random questions added"})
}

func (h *Handler) RemoveQuestion(c *gin.Context) {
	examID := c.Query("exam_id")
	eqID := c.Param("id")
	if examID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "exam_id query param required"})
		return
	}
	if err := h.UC.RemoveQuestion(c, examID, eqID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) GetExamMarks(c *gin.Context) {
	id := c.Param("id")
	out, err := h.UC.GetExamMarks(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) DownloadExamMarks(c *gin.Context) {
	id := c.Param("id")
	file, fileName, err := h.UC.DownloadExamMarks(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	
	if err := file.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write excel file"})
		return
	}
}
