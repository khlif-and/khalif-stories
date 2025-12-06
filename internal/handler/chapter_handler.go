package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"khalif-stories/internal/domain"
	"khalif-stories/pkg/utils"

)

type ChapterHandler struct {
	uc domain.ChapterUseCase
}

func NewChapterHandler(uc domain.ChapterUseCase) *ChapterHandler {
	return &ChapterHandler{uc: uc}
}

type CreateChapterRequest struct {
	StoryUUID string `form:"story_id" binding:"required"`
	// Field Title DIHAPUS
}

type AddChapterSlideRequest struct {
	Content  string `form:"content" binding:"required"`
	Sequence int    `form:"sequence" binding:"required"`
}

// CreateChapter godoc
// @Summary      Create a new chapter
// @Description  Create a chapter for a story (Title inherited from Story)
// @Tags         chapters
// @Accept       multipart/form-data
// @Produce      json
// @Param        story_id  formData  string  true  "Story UUID"
// @Success      201  {object}  domain.Chapter
// @Failure      500  {object}  utils.APIResponse
// @Router       /admin/chapters [post]
// @Security     BearerAuth
func (h *ChapterHandler) Create(c *gin.Context) {
	var req CreateChapterRequest
	if err := c.ShouldBind(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Hanya kirim StoryUUID
	res, err := h.uc.Create(c.Request.Context(), req.StoryUUID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, res)
}

// GetChapter godoc
// @Summary      Get chapter detail
// @Description  Get chapter by UUID with all slides
// @Tags         chapters
// @Produce      json
// @Param        uuid   path      string  true  "Chapter UUID"
// @Success      200  {object}  domain.Chapter
// @Failure      404  {object}  utils.APIResponse
// @Router       /chapters/{uuid} [get]
func (h *ChapterHandler) GetOne(c *gin.Context) {
	res, err := h.uc.GetByUUID(c.Request.Context(), c.Param("uuid"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, res)
}

// AddSlideToChapter godoc
// @Summary      Add slide to chapter
// @Description  Add slide (content, image, sound) to a chapter. Max 20 slides.
// @Tags         chapters
// @Accept       multipart/form-data
// @Produce      json
// @Param        uuid     path      string  true  "Chapter UUID"
// @Param        content  formData  string  true  "Content Text"
// @Param        sequence formData  int     true  "Sequence Number"
// @Param        image    formData  file    false "Slide Image"
// @Param        sound    formData  file    false "Slide Audio"
// @Success      201  {object}  domain.Slide
// @Failure      500  {object}  utils.APIResponse
// @Router       /admin/chapters/{uuid}/slides [post]
// @Security     BearerAuth
func (h *ChapterHandler) AddSlide(c *gin.Context) {
	var req AddChapterSlideRequest
	if err := c.ShouldBind(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	chapterUUID := c.Param("uuid")
	
	imageFile, imageHeader, _ := c.Request.FormFile("image")
	soundFile, soundHeader, _ := c.Request.FormFile("sound")

	res, err := h.uc.AddSlide(c.Request.Context(), chapterUUID, req.Content, req.Sequence, imageFile, imageHeader, soundFile, soundHeader)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, res)
}

// DeleteChapter godoc
// @Summary      Delete chapter
// @Description  Delete chapter and all its slides/assets
// @Tags         chapters
// @Produce      json
// @Param        uuid   path      string  true  "Chapter UUID"
// @Success      200  {object}  utils.APIResponse
// @Router       /admin/chapters/{uuid} [delete]
// @Security     BearerAuth
func (h *ChapterHandler) Delete(c *gin.Context) {
	if err := h.uc.Delete(c.Request.Context(), c.Param("uuid")); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessMessage(c, http.StatusOK, "chapter deleted")
}