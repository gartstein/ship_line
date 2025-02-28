package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"ship_line/services"
)

type dummyRepoForHandler struct {
	packSizes []int
}

func (d *dummyRepoForHandler) GetPackSizes() ([]int, error) {
	return d.packSizes, nil
}

func (d *dummyRepoForHandler) SetPackSizes(sizes []int) error {
	d.packSizes = sizes
	return nil
}

func (d *dummyRepoForHandler) DeletePackSize(size int) error {
	return nil
}

func TestHandler_CalcHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &dummyRepoForHandler{packSizes: []int{250, 500, 1000, 2000, 5000}}
	ps := services.NewPackService(repo)
	handler := Handler{ps: ps}
	router := gin.Default()
	router.GET("/v1/calc", handler.CalcHandler)

	t.Run("Missing query parameter", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/calc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid query parameter", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/calc?items=abc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Valid query parameter", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/calc?items=501", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var result map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err)

		// Verify that itemsOrdered and totalItemsUsed are as expected.
		itemsOrdered, err := strconv.Atoi(fmt.Sprintf("%v", result["itemsOrdered"]))
		assert.NoError(t, err)
		assert.Equal(t, 501, itemsOrdered)
		totalItemsUsed, err := strconv.Atoi(fmt.Sprintf("%v", result["totalItemsUsed"]))
		assert.NoError(t, err)
		assert.Equal(t, 750, totalItemsUsed)

		// Verify packsUsed.
		packsUsed := result["packsUsed"].(map[string]interface{})
		// Expected: {"500": 1, "250": 1}
		assert.Equal(t, 1, int(packsUsed["500"].(float64)))
		assert.Equal(t, 1, int(packsUsed["250"].(float64)))
	})
}
