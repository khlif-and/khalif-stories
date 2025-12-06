package handler

import (
	"errors"
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

// --- DTOs ---

type CreateCategoryRequest struct {
	Name string `form:"name" binding:"required"`
}

type UpdateCategoryRequest struct {
	Name string `form:"name"`
}

// --- HANDLERS ---

// CreateCategory godoc
// @Summary      Create a new category
// @Description  Create a new category with an image
// @Tags         categories
// @Accept       multipart/form-data
// @Produce      json
// @Param        name   formData  string  true  "Category Name"
// @Param        image  formData  file    false "Category Image"
// @Success      201  {object}  domain.Category
// @Failure      400  {object}  utils.APIResponse
// @Failure      409  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /admin/categories [post]
// @Security     BearerAuth
func (h *CategoryHandler) Create(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	file, header, _ := c.Request.FormFile("image")
	
	res, err := h.useCase.Create(c.Request.Context(), req.Name, file, header)
	if err != nil {
		if errors.Is(err, domain.ErrConflict) {
			utils.ErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, res)
}

// UpdateCategory godoc
// @Summary      Update a category
// @Description  Update category details
// @Tags         categories
// @Accept       multipart/form-data
// @Produce      json
// @Param        id     path      string  true  "Category UUID"
// @Param        name   formData  string  false "Category Name"
// @Param        image  formData  file    false "Category Image"
// @Success      200  {object}  domain.Category
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /admin/categories/{id} [put]
// @Security     BearerAuth
func (h *CategoryHandler) Update(c *gin.Context) {
	var req UpdateCategoryRequest
	_ = c.ShouldBind(&req)

	uuid := c.Param("id")
	file, header, _ := c.Request.FormFile("image")

	res, err := h.useCase.Update(c.Request.Context(), uuid, req.Name, file, header)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, domain.ErrConflict) {
			utils.ErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, res)
}

// DeleteCategory godoc
// @Summary      Delete a category
// @Description  Delete a category by UUID
// @Tags         categories
// @Produce      json
// @Param        id   path      string  true  "Category UUID"
// @Success      200  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /admin/categories/{id} [delete]
// @Security     BearerAuth
func (h *CategoryHandler) Delete(c *gin.Context) {
	if err := h.useCase.Delete(c.Request.Context(), c.Param("id")); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessMessage(c, http.StatusOK, "deleted")
}

// GetAllCategories godoc
// @Summary      Get all categories
// @Description  Retrieve all categories
// @Tags         categories
// @Produce      json
// @Success      200  {array}   domain.Category
// @Failure      500  {object}  utils.APIResponse
// @Router       /categories [get]
func (h *CategoryHandler) GetAll(c *gin.Context) {
	res, err := h.useCase.GetAll(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, res)
}

// GetCategory godoc
// @Summary      Get category by ID
// @Description  Retrieve a single category
// @Tags         categories
// @Produce      json
// @Param        id   path      string  true  "Category ID"
// @Success      200  {object}  domain.Category
// @Failure      404  {object}  utils.APIResponse
// @Router       /categories/{id} [get]
func (h *CategoryHandler) GetOne(c *gin.Context) {
	res, err := h.useCase.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, res)
}

// SearchCategories godoc
// @Summary      Search categories
// @Description  Search categories by name
// @Tags         categories
// @Produce      json
// @Param        q    query     string  true  "Search Query"
// @Success      200  {array}   domain.Category
// @Failure      400  {object}  utils.APIResponse
// @Router       /search/categories [get]
func (h *CategoryHandler) Search(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "query required")
		return
	}
	res, err := h.useCase.Search(c.Request.Context(), q)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, res)
}