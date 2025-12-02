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
		utils.ErrorResponse(c, http.StatusBadRequest, "Category name is required")
		return
	}

	file, header, _ := c.Request.FormFile("image")

	category, err := h.useCase.CreateCategory(name, file, header)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, category)
}

func (h *StoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	name := c.PostForm("name")
	file, header, _ := c.Request.FormFile("image")

	if name == "" && file == nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "No data to update")
		return
	}

	category, err := h.useCase.UpdateCategory(id, name, file, header)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, category)
}

func (h *StoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	err := h.useCase.DeleteCategory(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessMessage(c, http.StatusOK, "Category and related stories deleted successfully")
}

func (h *StoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.useCase.GetAllCategories()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, categories)
}

func (h *StoryHandler) GetCategory(c *gin.Context) {
	id := c.Param("id")

	category, err := h.useCase.GetCategory(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Category not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, category)
}

func (h *StoryHandler) CreateStory(c *gin.Context) {
	title := c.PostForm("title")
	desc := c.PostForm("description")
	catIDStr := c.PostForm("category_id")

	catID, err := strconv.Atoi(catIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid category_id")
		return
	}

	file, header, err := c.Request.FormFile("thumbnail")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Thumbnail image is required")
		return
	}

	story, err := h.useCase.CreateStory(title, desc, uint(catID), file, header)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, story)
}

func (h *StoryHandler) GetAllStories(c *gin.Context) {
	pagination := utils.GeneratePaginationFromRequest(c)

	stories, err := h.useCase.GetAllStories(pagination.Page, pagination.Limit, pagination.Sort)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponseWithMeta(c, http.StatusOK, stories, pagination)
}

func (h *StoryHandler) DeleteStory(c *gin.Context) {
	id := c.Param("id")

	err := h.useCase.DeleteStory(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessMessage(c, http.StatusOK, "Story deleted successfully")
}