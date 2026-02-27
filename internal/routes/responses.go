package routes

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) && len(validationErrs) > 0 {
		return validationErrorMessage(validationErrs[0])
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

func validationErrorMessage(fe validator.FieldError) string {
	field := humanFieldName(fe.Field())
	tag := strings.ToLower(fe.Tag())

	switch tag {
	case "required":
		return field + " is required"
	case "email":
		return "Enter a valid email address"
	case "min":
		if field == "Password" || field == "New password" {
			return field + " must be at least " + fe.Param() + " characters"
		}
		return field + " must be at least " + fe.Param() + " characters long"
	case "max":
		return field + " must be at most " + fe.Param() + " characters long"
	case "len":
		return field + " must be exactly " + fe.Param() + " characters long"
	default:
		return "Invalid value for " + strings.ToLower(field)
	}
}

func humanFieldName(field string) string {
	switch strings.ToLower(field) {
	case "username":
		return "Username"
	case "email":
		return "Email"
	case "password":
		return "Password"
	case "oldpassword":
		return "Current password"
	case "newpassword":
		return "New password"
	case "code":
		return "Verification code"
	case "title":
		return "Title"
	case "description":
		return "Description"
	case "deadline":
		return "Deadline"
	default:
		return field
	}
}
