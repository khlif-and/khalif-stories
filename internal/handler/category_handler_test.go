package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"khalif-stories/internal/domain"
	"khalif-stories/internal/handler"
	"khalif-stories/internal/mocks"

)

func TestCategoryHandler_GetAll(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockUC := new(mocks.CategoryUseCaseMock)
		h := handler.NewCategoryHandler(mockUC)

		expectedCategories := []domain.Category{
			{Name: "Horror", UUID: "uuid-1"},
			{Name: "Comedy", UUID: "uuid-2"},
		}

		mockUC.On("GetAll", mock.Anything).Return(expectedCategories, nil)

		r := gin.Default()
		r.GET("/categories", h.GetAll)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/categories", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		
		// Sesuaikan assertion ini dengan format JSON response kamu
		// Contoh jika response dibungkus "data":
		// assert.NotNil(t, response["data"]) 
	})

	t.Run("empty list", func(t *testing.T) {
		mockUC := new(mocks.CategoryUseCaseMock)
		h := handler.NewCategoryHandler(mockUC)

		mockUC.On("GetAll", mock.Anything).Return([]domain.Category{}, nil)

		r := gin.Default()
		r.GET("/categories", h.GetAll)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/categories", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}