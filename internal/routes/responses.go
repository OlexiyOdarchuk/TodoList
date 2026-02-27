package routes

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Err string `json:"err" binding:"required"`
}

type SuccessResponce struct {
	Message string `json:"message" binding:"required"`
}

func writeError(c *gin.Context, status int, err error) {
	c.JSON(status, ErrorResponse{Err: presentableErrorMessage(err)})
}

func writeSuccess(c *gin.Context, status int, message string) {
	c.JSON(status, SuccessResponce{Message: message})
}

func writeOK(c *gin.Context, message string) {
	writeSuccess(c, http.StatusOK, message)
}

func presentableErrorMessage(err error) string {
	if err == nil {
		return "Unexpected server error"
	}

	msg := strings.TrimSpace(err.Error())
	lower := strings.ToLower(msg)

	switch {
	case strings.Contains(lower, "users_email_key"):
		return "Email is already in use"
	case strings.Contains(lower, "users_username_key"):
		return "Username is already in use"
	case strings.Contains(lower, "duplicate key value violates unique constraint"),
		strings.Contains(lower, "unique constraint"):
		return "This value is already in use"
	case strings.Contains(lower, "violates foreign key constraint"):
		return "Related record was not found"
	case strings.Contains(lower, "violates not-null constraint"):
		return "Required field is missing"
	case strings.Contains(lower, "invalid input syntax"):
		return "Invalid input format"
	case strings.Contains(lower, "pq:"),
		strings.Contains(lower, "sqlstate"),
		strings.Contains(lower, "sql:"):
		return "Database request failed. Please try again"
	default:
		return msg
	}
}
