package domain

import "time"

// Category represents an expense category
type Category struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// CategoryRepository defines the interface for category data operations
type CategoryRepository interface {
	Create(category *Category) error
	GetByID(id int) (*Category, error)
	GetAll() ([]*Category, error)
	Update(category *Category) error
	Delete(id int) error
}
