package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"khalif-stories/internal/domain"
	"khalif-stories/pkg/utils"

)

type PreferenceHandler struct {
	uc domain.PreferenceUseCase
}

func NewPreferenceHandler(uc domain.PreferenceUseCase) *PreferenceHandler {
	return &PreferenceHandler{uc: uc}
}

type SavePreferencesRequest struct {
	StoryCategories  []string `json:"story_categories"`
	DakwahCategories []string `json:"dakwah_categories"`
	HadistCategories []string `json:"hadist_categories"`
}

func (h *PreferenceHandler) Save(c *gin.Context) {
	var req SavePreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if len(req.StoryCategories) > 5 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Maksimal pilih 5 Kategori Story")
		return
	}

	if len(req.DakwahCategories) > 5 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Maksimal pilih 5 Kategori Dakwah")
		return
	}

	if len(req.HadistCategories) > 5 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Maksimal pilih 5 Kategori Hadist")
		return
	}

	userID := c.GetString("user_id")
	err := h.uc.SavePreferences(c.Request.Context(), userID, req.StoryCategories, req.DakwahCategories, req.HadistCategories)
	
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessMessage(c, http.StatusOK, "preferences saved")
}