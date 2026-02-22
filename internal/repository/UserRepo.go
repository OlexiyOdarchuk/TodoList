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
	//TODO
	return &models.User{}, nil
}

func (r *UsersRepository) GetByEmail(email string) (*models.User, error) {
	//TODO
	return &models.User{}, nil
}

func (r *UsersRepository) GetByUsername(username string) (*models.User, error) {
	//TODO
	return &models.User{}, nil
}

func (r *UsersRepository) Create(u *models.User) error {
	//TODO
	return nil
}

func (r UsersRepository) Update(u *models.User) error {
	//TODO
	return nil
}

func (r *UsersRepository) Delete(id string) error {
	//TODO
	return nil
}

func (r *UsersRepository) GetByOAuth(provider, oauthID string) (*models.User, error) {
	//TODO
	return &models.User{}, nil
}
