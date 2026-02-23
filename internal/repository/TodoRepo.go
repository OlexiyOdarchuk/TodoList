package repository

import (
	"errors"
	"todolist/internal/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type TodoRepository struct {
	db *sqlx.DB
}

func NewTodoRepository(db *sqlx.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) Create(todo *models.Todo) error {
	query := `
		INSERT INTO todos (id, user_id, title, description, completed, created_at, updated_at, deadline)
		VALUES (:id, :user_id, :title, :description, :completed, :created_at, :updated_at, :deadline)
	`
	_, err := r.db.NamedExec(query, todo)
	return err
}

func (r *TodoRepository) GetListByUserID(userId string) ([]models.Todo, error) {
	query := `
	SELECT * FROM todos WHERE user_id = $1
	`
	var todos []models.Todo
	err := r.db.Select(&todos, query, userId)
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (r *TodoRepository) GetByID(id, userId string) (*models.Todo, error) {
	query := `SELECT * FROM todos WHERE id = $1 AND user_id = $2`
	var todo models.Todo
	err := r.db.Get(&todo, query, id, userId)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *TodoRepository) Update(todo *models.Todo) error {
	_, err := r.db.NamedExec(`
		UPDATE todos 
		SET title = :title, description = :description, completed = :completed, updated_at = CURRENT_TIMESTAMP, deadline = :deadline
		WHERE id = :id`, todo)

	if err != nil {
		return err
	}

	return nil
}

func (r *TodoRepository) Delete(id, userId string) error {
	res, err := r.db.Exec("DELETE FROM todos WHERE id = $1 AND user_id = $2", id, userId)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
