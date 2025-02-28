package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"ship_line/services"
)

// dummyRepo implements services.PackRepository for testing.
type dummyRepo struct {
	packSizes []int
}

func (d *dummyRepo) GetPackSizes() ([]int, error) {
	return d.packSizes, nil
}

func (d *dummyRepo) SetPackSizes(sizes []int) error {
	d.packSizes = sizes
	return nil
}

func (d *dummyRepo) DeletePackSize(size int) error {
	return nil
}

func TestHandler_UpdatePackSizes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Initialize dummy repository with initial pack sizes.
	repo := &dummyRepo{packSizes: []int{250, 500, 1000}}
	ps := services.NewPackService(repo)
	handler := &Handler{ps: ps}

	// Create a Gin router and register the endpoint.
	router := gin.Default()
	router.PUT("/v1/pack-sizes", handler.UpdatePackSizes)

	t.Run("Invalid JSON", func(t *testing.T) {
		req, err := http.NewRequest("PUT", "/v1/pack-sizes", bytes.NewBufferString("invalid"))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var resp map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "invalid JSON payload", resp["error"])
	})

	t.Run("Empty pack_sizes", func(t *testing.T) {
		payload := PackSizesPayload{PackSizes: []int{}}
		data, err := json.Marshal(payload)
		require.NoError(t, err)
		req, err := http.NewRequest("PUT", "/v1/pack-sizes", bytes.NewBuffer(data))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var resp map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "pack_sizes cannot be empty", resp["error"])
	})

	t.Run("Negative pack size", func(t *testing.T) {
		payload := PackSizesPayload{PackSizes: []int{250, -100}}
		data, _ := json.Marshal(payload)
		req, _ := http.NewRequest("PUT", "/v1/pack-sizes", bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var resp map[string]string
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "pack_sizes must be positive integers", resp["error"])
	})
	t.Run("Zero pack size", func(t *testing.T) {
		payload := PackSizesPayload{PackSizes: []int{250, 0, 100}}
		data, err := json.Marshal(payload)
		require.NoError(t, err)
		req, err := http.NewRequest("PUT", "/v1/pack-sizes", bytes.NewBuffer(data))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var resp map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "pack_sizes must be positive integers", resp["error"])
	})

	t.Run("Valid update", func(t *testing.T) {
		newSizes := []int{300, 600, 1200}
		payload := PackSizesPayload{PackSizes: newSizes}
		data, err := json.Marshal(payload)
		require.NoError(t, err)
		req, err := http.NewRequest("PUT", "/v1/pack-sizes", bytes.NewBuffer(data))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "success", resp["status"])

		// Because numbers in JSON are decoded as float64, convert before comparing.
		var updatedSizes []int
		for _, v := range resp["pack_sizes"].([]interface{}) {
			updatedSizes = append(updatedSizes, int(v.(float64)))
		}
		assert.Equal(t, newSizes, updatedSizes)

		// Verify that the dummy repository was updated.
		sizes, err := repo.GetPackSizes()
		require.NoError(t, err)
		assert.Equal(t, newSizes, sizes)
	})
}
