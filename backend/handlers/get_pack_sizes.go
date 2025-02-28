package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPackSizes is a Gin handler for retrieving the configured pack sizes.
// It responds with JSON in the form: { "pack_sizes": [250, 500, 1000, ...] }
func (h *Handler) GetPackSizes(c *gin.Context) {
	sizes, err := h.ps.GetPackSizes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sizes)
}
