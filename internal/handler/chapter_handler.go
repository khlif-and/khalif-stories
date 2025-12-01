package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"khalif-stories/internal/domain"

)

type ChapterHandler struct {
	useCase domain.StoryUseCase
}

func NewChapterHandler(u domain.StoryUseCase) *ChapterHandler {
	return &ChapterHandler{useCase: u}
}

func (h *ChapterHandler) AddSlide(c *gin.Context) {
	storyIDStr := c.Param("id")
	storyID, err := strconv.Atoi(storyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid story ID"})
		return
	}

	content := c.PostForm("content")
	seqStr := c.PostForm("sequence")
	sequence, _ := strconv.Atoi(seqStr)

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Slide image is required"})
		return
	}

	slide, err := h.useCase.AddSlide(uint(storyID), content, sequence, file, header)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": slide})
}