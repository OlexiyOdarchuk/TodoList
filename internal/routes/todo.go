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

// Create godoc
// @Summary Create a new todo
// @Description Create a new todo for the authenticated user
// @Tags todos
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param input body CreateTodoInput true "Todo payload"
// @Success 200 {object} models.Todo
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/todos [post]
func (h *TodoHandler) Create(c *gin.Context) {
	userID := c.MustGet("user").(string)

	var input CreateTodoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, err)
		return
	}

	todo, err := h.ts.CreateTodo(userID, input.Title, input.Description, input.Deadline)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, todo)
}

// GetAll godoc
// @Summary Get all todos
// @Description Get all todos for the authenticated user
// @Tags todos
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} models.Todo
// @Failure 500 {object} ErrorResponse
// @Router /api/todos [get]
func (h *TodoHandler) GetAll(c *gin.Context) {
	userID := c.MustGet("user").(string)

	todos, err := h.ts.GetTodosByUserID(userID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, todos)
}

// Delete godoc
// @Summary Delete a todo
// @Description Delete a todo by ID
// @Tags todos
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Todo ID"
// @Success 200 {object} SuccessResponce
// @Failure 500 {object} ErrorResponse
// @Router /api/todos/{id} [delete]
func (h *TodoHandler) Delete(c *gin.Context) {
	userID := c.MustGet("user").(string)
	todoID := c.Param("id")

	err := h.ts.DeleteTodo(userID, todoID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, err)
		return
	}

	writeOK(c, "Todo deleted successfully")
}

// Update godoc
// @Summary Update a todo
// @Description Update an existing todo
// @Tags todos
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Todo ID"
// @Param input body models.Todo true "Todo object"
// @Success 200 {object} SuccessResponce
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/todos/{id} [put]
func (h *TodoHandler) Update(c *gin.Context) {
	userID := c.MustGet("user").(string)
	todoID := c.Param("id")

	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		writeError(c, http.StatusBadRequest, err)
		return
	}

	todo.Id = todoID
	if err := h.ts.UpdateTodo(userID, &todo); err != nil {
		writeError(c, http.StatusInternalServerError, err)
		return
	}

	writeOK(c, "Todo updated successfully")
}
