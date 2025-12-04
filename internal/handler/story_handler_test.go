package handler_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
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

func TestStoryHandler_GetAll(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockUC := new(mocks.StoryUseCaseMock)
		h := handler.NewStoryHandler(mockUC)

		expectedStories := []domain.Story{
			{Title: "Story 1", UUID: "uuid-1"},
			{Title: "Story 2", UUID: "uuid-2"},
		}

		mockUC.On("GetAll", mock.Anything, 1, 10, "desc").Return(expectedStories, nil)

		r := gin.Default()
		r.GET("/stories", h.GetAll)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/stories?page=1&limit=10&sort=desc", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
	})
}

func TestStoryHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockUC := new(mocks.StoryUseCaseMock)
		h := handler.NewStoryHandler(mockUC)

		mockUC.On("Create", mock.Anything, "Title", "Desc", uint(1), mock.Anything, mock.Anything).
			Return(&domain.Story{Title: "Title"}, nil)

		r := gin.Default()
		r.POST("/stories", h.Create)

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		_ = writer.WriteField("title", "Title")
		_ = writer.WriteField("description", "Desc")
		_ = writer.WriteField("category_id", "1")

		part, _ := writer.CreateFormFile("file", "test.jpg")
		part.Write([]byte("dummy image content"))

		writer.Close()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/stories", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestStoryHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockUC := new(mocks.StoryUseCaseMock)
		h := handler.NewStoryHandler(mockUC)

		mockUC.On("Delete", mock.Anything, "uuid-123").Return(nil)

		r := gin.Default()
		r.DELETE("/stories/:id", h.Delete)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/stories/uuid-123", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}