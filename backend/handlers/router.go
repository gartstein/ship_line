// Package handlers defines HTTP route handlers for the Ship Line API.
// It configures the Gin router and maps endpoints to corresponding service methods.
package handlers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"ship_line/services"
	"time"
)

// Handler embeds the Gin engine and holds a reference to the PackService.
// It is used to handle HTTP requests for pack operations.
type Handler struct {
	*gin.Engine
	ps *services.PackService
}

// SetupRouter creates and configures the Gin router, setting up CORS policies and route handlers.
// It returns a Handler that encapsulates the router and associated pack service.
func SetupRouter(ps *services.PackService) *Handler {
	router := gin.Default()
	// Configure CORS to allow requests from specific origins and methods.
	router.Use(cors.New(cors.Config{
		// TODO: move to config
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodOptions, http.MethodPost, http.MethodDelete},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	handler := &Handler{
		Engine: router,
		ps:     ps,
	}

	// Define the route for pack calculations.
	router.GET("/v1/calc", handler.CalcHandler)
	// Define the route to update pack sizes.
	router.PUT("/v1/pack-sizes", handler.UpdatePackSizes)
	// Define the route to retrieve pack sizes.
	router.GET("/v1/pack-sizes", handler.GetPackSizes)
	// Define the route to delete a specific pack size.
	router.DELETE("/v1/pack-sizes/:size", handler.DeletePackSizeHandler)

	return handler
}
