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

func (h *StoryHandler) Create(c *gin.Context) {
	catID, _ := strconv.Atoi(c.PostForm("category_id"))
	file, header, err := c.Request.FormFile("thumbnail")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "thumbnail required")
		return
	}
	res, err := h.useCase.Create(c.PostForm("title"), c.PostForm("description"), uint(catID), file, header)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, res)
}

func (h *StoryHandler) GetAll(c *gin.Context) {
	p := utils.GeneratePaginationFromRequest(c)
	res, err := h.useCase.GetAll(p.Page, p.Limit, p.Sort)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponseWithMeta(c, http.StatusOK, res, p)
}

func (h *StoryHandler) Search(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "query required")
		return
	}
	res, err := h.useCase.Search(q)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, res)
}

func (h *StoryHandler) Delete(c *gin.Context) {
	if err := h.useCase.Delete(c.Param("id")); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessMessage(c, http.StatusOK, "deleted")
}

func (h *StoryHandler) AddSlide(c *gin.Context) {
	seq, _ := strconv.Atoi(c.PostForm("sequence"))
	file, header, _ := c.Request.FormFile("image")
	res, err := h.useCase.AddSlide(c.Param("id"), c.PostForm("content"), seq, file, header)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, res)
}