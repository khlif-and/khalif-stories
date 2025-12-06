package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"khalif-stories/internal/domain"
	"khalif-stories/pkg/utils"

)

type StoryHandler struct {
	uc domain.StoryUseCase
}

func NewStoryHandler(uc domain.StoryUseCase) *StoryHandler {
	return &StoryHandler{uc: uc}
}

type CreateStoryRequest struct {
	Title       string `form:"title" binding:"required"`
	Description string `form:"description" binding:"required"`
	CategoryID  string `form:"category_id" binding:"required"`
}

type UpdateStoryRequest struct {
	Title       string `form:"title"`
	Description string `form:"description"`
	CategoryID  string `form:"category_id"`
	Status      string `form:"status"`
}

type AddSlideRequest struct {
	Content  string `form:"content" binding:"required"`
	Sequence int    `form:"sequence" binding:"required"`
}

// CreateStory godoc
// @Summary      Create a new story
// @Description  Create a new story with thumbnail
// @Tags         stories
// @Accept       multipart/form-data
// @Produce      json
// @Param        title        formData  string  true  "Title"
// @Param        description  formData  string  true  "Description"
// @Param        category_id  formData  string  true  "Category UUID"
// @Param        file         formData  file    true  "Thumbnail Image"
// @Success      201  {object}  domain.Story
// @Failure      400  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /admin/stories [post]
// @Security     BearerAuth
func (h *StoryHandler) Create(c *gin.Context) {
	var req CreateStoryRequest
	if err := c.ShouldBind(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "thumbnail image is required")
		return
	}

	userID := c.GetString("user_id")
	story, err := h.uc.Create(c.Request.Context(), req.Title, req.Description, req.CategoryID, userID, file, header)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, story)
}

// UpdateStory godoc
// @Summary      Update a story
// @Description  Update story details
// @Tags         stories
// @Accept       multipart/form-data
// @Produce      json
// @Param        uuid         path      string  true  "Story UUID"
// @Param        title        formData  string  false "Title"
// @Param        description  formData  string  false "Description"
// @Param        category_id  formData  string  false "Category UUID"
// @Param        status       formData  string  false "Status"
// @Param        file         formData  file    false "Thumbnail Image"
// @Success      200  {object}  domain.Story
// @Failure      500  {object}  utils.APIResponse
// @Router       /admin/stories/{uuid} [put]
// @Security     BearerAuth
func (h *StoryHandler) Update(c *gin.Context) {
	var req UpdateStoryRequest
	_ = c.ShouldBind(&req)

	uuid := c.Param("uuid")
	file, header, _ := c.Request.FormFile("file")

	story, err := h.uc.Update(c.Request.Context(), uuid, req.Title, req.Description, req.CategoryID, req.Status, file, header)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, story)
}

// GetAllStories godoc
// @Summary      Get all stories
// @Description  Get stories with pagination
// @Tags         stories
// @Produce      json
// @Param        page   query     int     false "Page"
// @Param        limit  query     int     false "Limit"
// @Param        sort   query     string  false "Sort (e.g., created_at desc)"
// @Success      200  {array}   domain.Story
// @Failure      500  {object}  utils.APIResponse
// @Router       /stories [get]
func (h *StoryHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sort := c.DefaultQuery("sort", "created_at desc")

	stories, err := h.uc.GetAll(c.Request.Context(), page, limit, sort)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, stories)
}

// GetStory godoc
// @Summary      Get story by UUID
// @Description  Retrieve a single story
// @Tags         stories
// @Produce      json
// @Param        uuid   path      string  true  "Story UUID"
// @Success      200  {object}  domain.Story
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /stories/{uuid} [get]
func (h *StoryHandler) GetOne(c *gin.Context) {
	uuid := c.Param("uuid")
	story, err := h.uc.GetByUUID(c.Request.Context(), uuid)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "story not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, story)
}

// SearchStories godoc
// @Summary      Search stories
// @Description  Search stories by title or description
// @Tags         stories
// @Produce      json
// @Param        q    query     string  true  "Search Query"
// @Success      200  {array}   domain.Story
// @Failure      500  {object}  utils.APIResponse
// @Router       /search/stories [get]
func (h *StoryHandler) Search(c *gin.Context) {
	q := c.Query("q")
	// PERBAIKAN: Menambahkan validasi wajib input q agar mirip Category
	if q == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "query required")
		return
	}
	stories, err := h.uc.Search(c.Request.Context(), q)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, stories)
}

// DeleteStory godoc
// @Summary      Delete a story
// @Description  Delete a story by UUID
// @Tags         stories
// @Produce      json
// @Param        uuid path      string  true  "Story UUID"
// @Success      200  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /admin/stories/{uuid} [delete]
// @Security     BearerAuth
func (h *StoryHandler) Delete(c *gin.Context) {
	uuid := c.Param("uuid")
	if err := h.uc.Delete(c.Request.Context(), uuid); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessMessage(c, http.StatusOK, "story deleted successfully")
}

// AddSlide godoc
// @Summary      Add a slide to story
// @Description  Add content slide to story
// @Tags         stories
// @Accept       multipart/form-data
// @Produce      json
// @Param        uuid     path      string  true  "Story UUID"
// @Param        content  formData  string  true  "Content Text"
// @Param        sequence formData  int     true  "Sequence Number"
// @Param        file     formData  file    false "Slide Image"
// @Success      201  {object}  domain.Slide
// @Failure      500  {object}  utils.APIResponse
// @Router       /admin/stories/{uuid}/slides [post]
// @Security     BearerAuth
func (h *StoryHandler) AddSlide(c *gin.Context) {
	var req AddSlideRequest
	if err := c.ShouldBind(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	storyUUID := c.Param("uuid")
	file, header, _ := c.Request.FormFile("file")

	slide, err := h.uc.AddSlide(c.Request.Context(), storyUUID, req.Content, req.Sequence, file, header)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, slide)
}