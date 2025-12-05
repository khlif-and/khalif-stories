package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"khalif-stories/internal/domain"

)

type StoryHandler struct {
	uc domain.StoryUseCase
}

func NewStoryHandler(uc domain.StoryUseCase) *StoryHandler {
	return &StoryHandler{uc: uc}
}

func (h *StoryHandler) Create(c *gin.Context) {
	title := c.PostForm("title")
	desc := c.PostForm("description")
	categoryID := c.PostForm("category_id")
	userID := c.GetString("user_id") 

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Thumbnail image is required"})
		return
	}

	story, err := h.uc.Create(c.Request.Context(), title, desc, categoryID, userID, file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Story created successfully",
		"data":    story,
	})
}

func (h *StoryHandler) Update(c *gin.Context) {
	uuid := c.Param("uuid")
	title := c.PostForm("title")
	desc := c.PostForm("description")
	categoryID := c.PostForm("category_id")
	status := c.PostForm("status")

	file, header, _ := c.Request.FormFile("file")

	story, err := h.uc.Update(c.Request.Context(), uuid, title, desc, categoryID, status, file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Story updated successfully",
		"data":    story,
	})
}

func (h *StoryHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sort := c.DefaultQuery("sort", "created_at desc")

	stories, err := h.uc.GetAll(c.Request.Context(), page, limit, sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stories})
}

func (h *StoryHandler) Search(c *gin.Context) {
	query := c.Query("q")
	stories, err := h.uc.Search(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stories})
}

func (h *StoryHandler) Delete(c *gin.Context) {
	uuid := c.Param("uuid")
	if err := h.uc.Delete(c.Request.Context(), uuid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Story deleted successfully"})
}

func (h *StoryHandler) AddSlide(c *gin.Context) {
	storyUUID := c.Param("uuid")
	content := c.PostForm("content")
	sequence, _ := strconv.Atoi(c.PostForm("sequence"))

	file, header, _ := c.Request.FormFile("file")

	slide, err := h.uc.AddSlide(c.Request.Context(), storyUUID, content, sequence, file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Slide added successfully",
		"data":    slide,
	})
}