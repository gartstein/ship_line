package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"ship_line/swagger"
)

// MockControllerForGet simulates the controller returning PackSizesPayload.
type MockControllerForGet struct {
	response *swagger.PackSizesPayload
	err      error
}

func (m *MockControllerForGet) GetPackSizesHandler(c *gin.Context) {
	if m.err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": m.err.Error()})
		return
	}
	c.JSON(http.StatusOK, m.response)
}

func TestHandler_GetPackSizes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("ValidGetPackSizes", func(t *testing.T) {
		mockController := &MockControllerForGet{
			response: &swagger.PackSizesPayload{
				PackSizes: []int{250, 500, 1000},
			},
		}

		router := gin.Default()
		router.GET("/v1/pack-sizes", mockController.GetPackSizesHandler)

		req, _ := http.NewRequest("GET", "/v1/pack-sizes", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp swagger.PackSizesPayload
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		expectedResponse := swagger.PackSizesPayload{
			PackSizes: []int{250, 500, 1000},
		}

		assert.Equal(t, expectedResponse, resp)
	})

	t.Run("EmptyPackSizes", func(t *testing.T) {
		mockController := &MockControllerForGet{
			response: &swagger.PackSizesPayload{
				PackSizes: []int{}, // Empty list
			},
		}

		router := gin.Default()
		router.GET("/v1/pack-sizes", mockController.GetPackSizesHandler)

		req, _ := http.NewRequest("GET", "/v1/pack-sizes", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp swagger.PackSizesPayload
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		assert.Empty(t, resp.PackSizes)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockController := &MockControllerForGet{err: assert.AnError}

		router := gin.Default()
		router.GET("/v1/pack-sizes", mockController.GetPackSizesHandler)

		req, _ := http.NewRequest("GET", "/v1/pack-sizes", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})
}
