package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// DeletePackSizeHandler handles DELETE /pack-sizes/{size}
func (h *Handler) DeletePackSizeHandler(c *gin.Context) {
	size, err := strconv.Atoi(c.Param("size"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid size"})
		return
	}

	err = h.ps.DeletePackSizeHandler(size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
