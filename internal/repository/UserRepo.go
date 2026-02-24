package repository

import (
	"todolist/internal/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type UsersRepository struct {
	db *sqlx.DB
}

func NewUsersRepository(db *sqlx.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) GetByID(id string) (*models.User, error) {
	var u models.User
	err := r.db.Get(&u, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UsersRepository) GetByEmail(email string) (*models.User, error) {
	var u models.User
	err := r.db.Get(&u, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UsersRepository) GetByUsername(username string) (*models.User, error) {
	var u models.User
	err := r.db.Get(&u, "SELECT * FROM users WHERE username = $1", username)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UsersRepository) CreateWithDefaultTodo(u *models.User, todo *models.Todo) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	_, err = tx.NamedExec(`INSERT INTO users (id, username, email, password_hash, is_verified, verification_code, verification_code_expires, oauth_provider, oauth_id) VALUES (:id, :username, :email, :password_hash, :is_verified, :verification_code, :verification_code_expires, :oauth_provider, :oauth_id)`, u)
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(`
		INSERT INTO todos (id, user_id, title, description, completed, created_at, updated_at, deadline)
		VALUES (:id, :user_id, :title, :description, :completed, :created_at, :updated_at, :deadline)
	`, todo)
	if err != nil {
		return err
	}

	return nil
}

func (r UsersRepository) Update(u *models.User) error {
	query := `UPDATE users 
			SET username = :username, email = :email, pending_email = :pending_email, password_hash = :password_hash, is_verified = :is_verified, verification_code = :verification_code, verification_code_expires = :verification_code_expires
			WHERE id = :id`
	_, err := r.db.NamedExec(query, &u)
	return err
}

func (r *UsersRepository) Delete(id string) error {
	_, err := r.db.NamedExec(`DELETE FROM users WHERE id = :id`, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *UsersRepository) GetByOAuth(provider, oauthID string) (*models.User, error) {
	query := `
	SELECT * FROM users WHERE oauth_id = $1 AND provider = $2
	`
	var user models.User
	err := r.db.Get(&user, query, oauthID, provider)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
