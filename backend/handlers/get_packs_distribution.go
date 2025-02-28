package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CalcHandler handles GET /v1/calc?items=X.
func (h *Handler) CalcHandler(c *gin.Context) {
	itemsStr := c.Query("items")
	if itemsStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'items' query param"})
		return
	}
	items, err := strconv.Atoi(itemsStr)
	if err != nil || items < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'items' value"})
		return
	}

	// TODO: get from config
	const maxOrder = 1000000000000
	if items > maxOrder {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("order %d exceeds maximum allowed value of %d", items, maxOrder)})
		return
	}

	result, err := h.ps.CalculatePacks(items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
