package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck godoc
// @Summary Health Check
// @Description Check if the server is running
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func Register(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}
