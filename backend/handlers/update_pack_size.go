package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PackSizesPayload defines the expected JSON body for updating pack sizes.
type PackSizesPayload struct {
	PackSizes []int `json:"pack_sizes"`
}

// UpdatePackSizes is a Gin handler for updating pack sizes.
// It expects a JSON payload like: { "pack_sizes": [250,500,1000] }
func (h *Handler) UpdatePackSizes(c *gin.Context) {
	var payload PackSizesPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON payload"})
		return
	}
	if len(payload.PackSizes) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pack_sizes cannot be empty"})
		return
	}
	for _, size := range payload.PackSizes {
		if size <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "pack_sizes must be positive integers"})
			return
		}
	}

	err := h.ps.UpdatePackSizes(payload.PackSizes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "pack_sizes": payload.PackSizes})
}
