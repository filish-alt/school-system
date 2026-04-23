package teacher

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	teacherdto "school-exam/internal/dto/teacher"
	"school-exam/internal/module/teacher"
)

type Handler struct {
	UC *teacher.Usecase
}

func New(uc *teacher.Usecase) *Handler {
	return &Handler{UC: uc}
}

func (h *Handler) CreateQuestionBank(c *gin.Context) {
	var req teacherdto.CreateQuestionBankRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	id, err := h.UC.CreateQuestionBank(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) ListQuestionBanks(c *gin.Context) {
	var q teacherdto.ListQuery
	_ = c.ShouldBindQuery(&q)
	out, err := h.UC.ListQuestionBanks(c, q.Page, q.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) CreateQuestion(c *gin.Context) {
	var req teacherdto.CreateQuestionRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	id, err := h.UC.CreateQuestion(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) ListQuestions(c *gin.Context) {
	bankID := c.Query("bank_id")
	if bankID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bank_id required"})
		return
	}
	var q teacherdto.ListQuery
	_ = c.ShouldBindQuery(&q)
	out, err := h.UC.ListQuestions(c, bankID, q.Page, q.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) UpdateQuestion(c *gin.Context) {
	var req teacherdto.UpdateQuestionRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	if err := h.UC.UpdateQuestion(c, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) DeleteQuestion(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}
	if err := h.UC.DeleteQuestion(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) CreateOption(c *gin.Context) {
	var req teacherdto.CreateOptionRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	id, err := h.UC.CreateOption(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) ListOptions(c *gin.Context) {
	qid := c.Query("question_id")
	if qid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "question_id required"})
		return
	}
	var q teacherdto.ListQuery
	_ = c.ShouldBindQuery(&q)
	out, err := h.UC.ListOptions(c, qid, q.Page, q.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) UpdateOption(c *gin.Context) {
	var req teacherdto.UpdateOptionRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	if err := h.UC.UpdateOption(c, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) DeleteOption(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}
	if err := h.UC.DeleteOption(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) ListMyStudents(c *gin.Context) {
	var q teacherdto.ListQuery
	_ = c.ShouldBindQuery(&q)
	out, err := h.UC.ListMyStudents(c, q.Page, q.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) ListMyAssignments(c *gin.Context) {
	var q teacherdto.ListQuery
	_ = c.ShouldBindQuery(&q)
	out, err := h.UC.ListMyAssignments(c, q.Page, q.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) ImportQuestions(c *gin.Context) {
	bankID := c.PostForm("bank_id")
	if bankID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bank_id required"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid csv format"})
		return
	}

	if err := h.UC.ImportQuestionsFromCSV(c, bankID, records); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "questions imported successfully"})
}

func (h *Handler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image required"})
		return
	}

	// Create uploads directory if it doesn't exist
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		err = os.Mkdir("uploads", 0755)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload directory"})
			return
		}
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := uuid.New().String() + ext
	savePath := filepath.Join("uploads", filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	// Return the relative URL
	url := fmt.Sprintf("/uploads/%s", filename)
	c.JSON(http.StatusOK, gin.H{"url": url})
}