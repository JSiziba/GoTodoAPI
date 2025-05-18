package repository

import (
	"errors"
	"todo/internal/models"

	"gorm.io/gorm"
)

// TodoRepository handles database operations for Todo
type TodoRepository struct {
	db *gorm.DB
}

// NewTodoRepository creates a new TodoRepository
func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

// Create adds a new Todo to the database
func (r *TodoRepository) Create(todo *models.Todo) error {
	return r.db.Create(todo).Error
}

// FindAll retrieves all todos from the database
func (r *TodoRepository) FindAll() ([]models.Todo, error) {
	var todos []models.Todo
	err := r.db.Find(&todos).Error
	return todos, err
}

// FindByID retrieves a todo by its ID
func (r *TodoRepository) FindByID(id uint) (*models.Todo, error) {
	var todo models.Todo
	err := r.db.First(&todo, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &todo, nil
}

// Update updates a todo in the database
func (r *TodoRepository) Update(todo *models.Todo) error {
	return r.db.Save(todo).Error
}

// Delete removes a todo from the database
func (r *TodoRepository) Delete(id uint) error {
	return r.db.Delete(&models.Todo{}, id).Error
}
