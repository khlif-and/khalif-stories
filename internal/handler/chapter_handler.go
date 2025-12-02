package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"khalif-stories/internal/domain"
	"khalif-stories/pkg/utils"

)

type ChapterHandler struct {
	useCase domain.StoryUseCase
}

func NewChapterHandler(u domain.StoryUseCase) *ChapterHandler {
	return &ChapterHandler{useCase: u}
}

func (h *ChapterHandler) AddSlide(c *gin.Context) {
	storyUUID := c.Param("id")

	content := c.PostForm("content")
	seqStr := c.PostForm("sequence")
	sequence, _ := strconv.Atoi(seqStr)

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Slide image is required")
		return
	}

	slide, err := h.useCase.AddSlide(storyUUID, content, sequence, file, header)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, slide)
}