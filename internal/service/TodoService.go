package service

import (
	"errors"
	"time"
	"todolist/internal/models"

	"github.com/google/uuid"
)

type TodoRepository interface {
	Create(todo *models.Todo) error
	GetListByUserID(userId string) ([]models.Todo, error)
	Update(todo *models.Todo) error
	Delete(id string, userId string) error
}

type TodoService struct {
	repo TodoRepository
}

func NewTodoService(repo TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) CreateTodo(userId, title, description string, deadline time.Time) (*models.Todo, error) {
	if userId == "" || title == "" {
		return nil, errors.New("userId or title is empty")
	}
	if deadline.Before(time.Now()) {
		return nil, errors.New("deadline is before now")
	}
	newTodo := &models.Todo{
		Id:          uuid.New().String(),
		UserID:      userId,
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Deadline:    deadline,
	}

	err := s.repo.Create(newTodo)
	if err != nil {
		return nil, err
	}
	return newTodo, nil
}

func (s *TodoService) UpdateTodo(userId string, newTodo *models.Todo) error {
	if newTodo.UserID != userId {
		return errors.New("wrong userId")
	}
	err := s.repo.Update(newTodo)
	if err != nil {
		return err
	}
	return nil
}

func (s *TodoService) DeleteTodo(userId string, todoId string) error {
	err := s.repo.Delete(todoId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (s *TodoService) GetTodosByUserID(userId string) ([]models.Todo, error) {
	return s.repo.GetListByUserID(userId)
}
