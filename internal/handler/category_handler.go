package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"khalif-stories/internal/domain"
	"khalif-stories/pkg/utils"

)

type CategoryHandler struct {
	useCase domain.CategoryUseCase
}

func NewCategoryHandler(u domain.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{useCase: u}
}

func (h *CategoryHandler) Create(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "name required")
		return
	}
	file, header, _ := c.Request.FormFile("image")
	res, err := h.useCase.Create(name, file, header)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, res)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	res, err := h.useCase.Update(c.Param("id"), c.PostForm("name"), nil, nil)
	file, header, _ := c.Request.FormFile("image")
	if file != nil {
		res, err = h.useCase.Update(c.Param("id"), c.PostForm("name"), file, header)
	}
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, res)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	if err := h.useCase.Delete(c.Param("id")); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessMessage(c, http.StatusOK, "deleted")
}

func (h *CategoryHandler) GetAll(c *gin.Context) {
	res, err := h.useCase.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, res)
}

func (h *CategoryHandler) GetOne(c *gin.Context) {
	res, err := h.useCase.Get(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "not found")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, res)
}

func (h *CategoryHandler) Search(c *gin.Context) {
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