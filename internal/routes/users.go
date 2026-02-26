package routes

import (
	"net/http"
	"todolist/internal/service"
	"todolist/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	us *service.UserService
}

func NewUserHandler(us *service.UserService) *UserHandler {
	return &UserHandler{us: us}
}

type RegisterInput struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=8"`
	Email    string `json:"email" binding:"required,email"`
}

type LoginInput struct {
	UserData string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type GoogleLoginInput struct {
	Token string `json:"token" binding:"required"`
}

type VerifyEmailRegisterInput struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

type UpdateUsernameInput struct {
	Username string `json:"username" binding:"required,min=3"`
}

type UpdatePasswordInput struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type RequestEmailUpdateInput struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyEmailInput struct {
	Code string `json:"code" binding:"required,len=6"`
}

type DeleteUserInput struct {
	Password string `json:"password"`
}

type UserProfileResponse struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	PendingEmail  string `json:"pending_email"`
	OauthProvider string `json:"oauth_provider"`
	HasPassword   bool   `json:"has_password"`
}

// Register godoc
// @Summary Register user
// @Description Creates a local user and sends email verification code
// @Tags auth
// @Accept json
// @Produce json
// @Param input body RegisterInput true "Registration payload"
// @Success 201 {object} service.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
func (uh *UserHandler) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, err)
		return
	}
	resp, err := uh.us.RegisterUser(input.Username, input.Email, input.Password)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// Login godoc
// @Summary Login user
// @Description Login by username/email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param input body LoginInput true "Login payload"
// @Success 200 {object} service.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/login [post]
func (uh *UserHandler) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, err)
		return
	}
	resp, err := uh.us.LoginUser(input.UserData, input.Password)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// VerifyEmail godoc
// @Summary Verify registration email
// @Description Verifies code from register flow and returns auth token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body VerifyEmailRegisterInput true "Email verification payload"
// @Success 200 {object} service.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/verify [post]
func (uh *UserHandler) VerifyEmail(c *gin.Context) {
	var input VerifyEmailRegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, err)
		return
	}
	resp, err := uh.us.VerifyEmail(input.Email, input.Code)
	if err != nil {
		writeError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GoogleLogin godoc
// @Summary Login with Google
// @Description Validates Google ID token and logs in or creates user
// @Tags auth
// @Accept json
// @Produce json
// @Param input body GoogleLoginInput true "Google login payload"
// @Success 200 {object} service.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/google [post]
func (uh *UserHandler) GoogleLogin(c *gin.Context) {
	var input GoogleLoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, err)
		return
	}

	oauthData, err := utils.ValidateGoogleToken(c.Request.Context(), input.Token)
	if err != nil {
		writeError(c, http.StatusBadRequest, err)
		return
	}

	resp, err := uh.us.LoginWithOAuth(c.Request.Context(), *oauthData)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetUser godoc
// @Summary Get current user
// @Description Returns profile of authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} UserProfileResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/me [get]
func (uh *UserHandler) GetUser(c *gin.Context) {
	userID := c.MustGet("user").(string)
	user, err := uh.us.GetUserByID(userID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, UserProfileResponse{
		ID:            user.Id,
		Username:      user.Username,
		Email:         user.Email,
		PendingEmail:  user.PendingEmail,
		OauthProvider: user.OauthProvider,
		HasPassword:   user.PasswordHash != "",
	})
}

// UpdateUsername godoc
// @Summary Update username
// @Description Updates username for authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param input body UpdateUsernameInput true "New username"
// @Success 200 {object} SuccessResponce
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/me [patch]
func (uh *UserHandler) UpdateUsername(c *gin.Context) {
	var input UpdateUsernameInput
	userID := c.MustGet("user").(string)
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, err)
		return
	}
	err := uh.us.UpdateUsername(userID, input.Username)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err)
		return
	}
	writeOK(c, "Username updated successfully")
}

// UpdatePassword godoc
// @Summary Update password
// @Description Updates password for authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param input body UpdatePasswordInput true "Password update payload"
// @Success 200 {object} SuccessResponce
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/me/password [put]
func (uh *UserHandler) UpdatePassword(c *gin.Context) {
	var input UpdatePasswordInput
	userID := c.MustGet("user").(string)
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, err)
		return
	}
	err := uh.us.UpdatePassword(userID, input.OldPassword, input.NewPassword)
	if err != nil {
		writeError(c, http.StatusBadRequest, err)
		return
	}
	writeOK(c, "Password updated successfully")
}

// RequestEmailUpdate godoc
// @Summary Request email update
// @Description Sends verification code to new email for authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param input body RequestEmailUpdateInput true "New email payload"
// @Success 200 {object} SuccessResponce
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/me/email [post]
func (uh *UserHandler) RequestEmailUpdate(c *gin.Context) {
	var input RequestEmailUpdateInput
	userID := c.MustGet("user").(string)
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, err)
		return
	}
	err := uh.us.RequestEmailUpdate(userID, input.Email)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err)
		return
	}
	writeOK(c, "Verification code sent to new email")
}

// VerifyEmailUpdate godoc
// @Summary Verify email update
// @Description Confirms code and updates email for authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param input body VerifyEmailInput true "Email verification code"
// @Success 200 {object} SuccessResponce
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/me/email [put]
func (uh *UserHandler) VerifyEmailUpdate(c *gin.Context) {
	var input VerifyEmailInput
	userID := c.MustGet("user").(string)
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, err)
		return
	}
	err := uh.us.VerifyEmailUpdate(userID, input.Code)
	if err != nil {
		writeError(c, http.StatusBadRequest, err)
		return
	}
	writeOK(c, "Email updated successfully")
}

// DeleteUser godoc
// @Summary Request user deletion
// @Description Sends verification code to current email before account deletion
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param input body DeleteUserInput true "Delete user payload"
// @Success 200 {object} SuccessResponce
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/me [delete]
func (uh *UserHandler) DeleteUser(c *gin.Context) {
	var input DeleteUserInput
	userID := c.MustGet("user").(string)
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, err)
		return
	}
	err := uh.us.DeleteUser(userID, input.Password)
	if err != nil {
		writeError(c, http.StatusBadRequest, err)
		return
	}
	writeOK(c, "Verification code sent to your email")
}

// VerifyEmailDelete godoc
// @Summary Confirm user deletion
// @Description Deletes authenticated user after verification code confirmation
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param input body VerifyEmailInput true "Delete verification code"
// @Success 200 {object} SuccessResponce
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/user/me/delete [put]
func (uh *UserHandler) VerifyEmailDelete(c *gin.Context) {
	var input VerifyEmailInput
	userID := c.MustGet("user").(string)
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, err)
		return
	}
	err := uh.us.VerifyEmailDelete(userID, input.Code)
	if err != nil {
		writeError(c, http.StatusBadRequest, err)
		return
	}
	writeOK(c, "User deleted successfully")
}
