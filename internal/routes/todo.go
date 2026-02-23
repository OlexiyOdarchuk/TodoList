package routes

import (
	"net/http"
	"time"
	"todolist/internal/models"
	"todolist/internal/service"

	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	ts *service.TodoService
}

func NewTodoHandler(ts *service.TodoService) *TodoHandler {
	return &TodoHandler{ts: ts}
}

type CreateTodoInput struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
}

func (h *TodoHandler) Create(c *gin.Context) {
	userID := c.MustGet("user").(string)

	var input CreateTodoInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo, err := h.ts.CreateTodo(userID, input.Title, input.Description, input.Deadline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) GetAll(c *gin.Context) {
	userID := c.MustGet("user").(string)

	todos, err := h.ts.GetTodosByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) Delete(c *gin.Context) {
	userID := c.MustGet("user").(string)
	todoID := c.Param("id")
	err := h.ts.DeleteTodo(userID, todoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *TodoHandler) Update(c *gin.Context) {

	userID := c.MustGet("user").(string)

	todoID := c.Param("id")

	var todo models.Todo

	if err := c.ShouldBindJSON(&todo); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		return

	}

	todo.Id = todoID

	if err := h.ts.UpdateTodo(userID, &todo); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})

	}

	c.Status(http.StatusOK)

}
