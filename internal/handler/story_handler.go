package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"khalif-stories/internal/domain"
	"khalif-stories/pkg/utils"

)

type StoryHandler struct {
	useCase domain.StoryUseCase
}

func NewStoryHandler(u domain.StoryUseCase) *StoryHandler {
	return &StoryHandler{useCase: u}
}

func (h *StoryHandler) CreateCategory(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category name is required"})
		return
	}

	file, header, _ := c.Request.FormFile("image")

	category, err := h.useCase.CreateCategory(name, file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": category})
}

func (h *StoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.useCase.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": categories})
}

func (h *StoryHandler) CreateStory(c *gin.Context) {
	title := c.PostForm("title")
	desc := c.PostForm("description")
	catIDStr := c.PostForm("category_id")

	catID, err := strconv.Atoi(catIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category_id"})
		return
	}

	file, header, err := c.Request.FormFile("thumbnail")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Thumbnail image is required"})
		return
	}

	story, err := h.useCase.CreateStory(title, desc, uint(catID), file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": story})
}

func (h *StoryHandler) GetAllStories(c *gin.Context) {
	pagination := utils.GeneratePaginationFromRequest(c)

	stories, err := h.useCase.GetAllStories(pagination.Page, pagination.Limit, pagination.Sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": stories,
		"meta": pagination,
	})
}

func (h *StoryHandler) DeleteStory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.useCase.DeleteStory(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Story deleted successfully"})
}