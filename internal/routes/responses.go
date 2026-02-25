package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Err string `json:"err" binding:"required"`
}

type SuccessResponce struct {
	Message string `json:"message" binding:"required"`
}

func writeError(c *gin.Context, status int, err error) {
	c.JSON(status, ErrorResponse{Err: err.Error()})
}

func writeSuccess(c *gin.Context, status int, message string) {
	c.JSON(status, SuccessResponce{Message: message})
}

func writeOK(c *gin.Context, message string) {
	writeSuccess(c, http.StatusOK, message)
}
