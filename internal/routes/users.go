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

type VerifyEmailInput struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

type UpdateUsernameInput struct {
	Username string `json:"username" binding:"required,min=3"`
}

type UpdatePasswordInput struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type RequestEmailUpdateInput struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyEmailUpdateInput struct {
	Code string `json:"code" binding:"required,len=6"`
}

func (uh *UserHandler) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := uh.us.RegisterUser(input.Username, input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (uh *UserHandler) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := uh.us.LoginUser(input.UserData, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (uh *UserHandler) VerifyEmail(c *gin.Context) {
	var input VerifyEmailInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := uh.us.VerifyEmail(input.Email, input.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (uh *UserHandler) GoogleLogin(c *gin.Context) {
	var input GoogleLoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	oauthData, err := utils.ValidateGoogleToken(c.Request.Context(), input.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := uh.us.LoginWithOAuth(c.Request.Context(), *oauthData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (uh *UserHandler) GetUser(c *gin.Context) {
	userId := c.MustGet("user").(string)
	user, err := uh.us.GetUserByID(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":             user.Id,
		"username":       user.Username,
		"email":          user.Email,
		"pending_email":  user.PendingEmail,
		"oauth_provider": user.OauthProvider,
		"has_password":   user.PasswordHash != "",
	})
}

func (uh *UserHandler) UpdateUsername(c *gin.Context) {
	var input UpdateUsernameInput
	userId := c.MustGet("user").(string)
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := uh.us.UpdateUsername(userId, input.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Username updated successfully"})
}

func (uh *UserHandler) UpdatePassword(c *gin.Context) {
	var input UpdatePasswordInput
	userId := c.MustGet("user").(string)
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := uh.us.UpdatePassword(userId, input.OldPassword, input.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

func (uh *UserHandler) RequestEmailUpdate(c *gin.Context) {
	var input RequestEmailUpdateInput
	userId := c.MustGet("user").(string)
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := uh.us.RequestEmailUpdate(userId, input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Verification code sent to new email"})
}

func (uh *UserHandler) VerifyEmailUpdate(c *gin.Context) {
	var input VerifyEmailUpdateInput
	userId := c.MustGet("user").(string)
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := uh.us.VerifyEmailUpdate(userId, input.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Email updated successfully"})
}
