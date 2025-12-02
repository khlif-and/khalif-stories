package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"khalif-stories/internal/domain"

)

type SearchHandler struct {
	repo domain.SearchRepository
}

func NewSearchHandler(repo domain.SearchRepository) *SearchHandler {
	return &SearchHandler{repo: repo}
}

func (h *SearchHandler) SearchCategories(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query param 'q' is required"})
		return
	}

	categories, err := h.repo.SearchCategories(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": categories})
}