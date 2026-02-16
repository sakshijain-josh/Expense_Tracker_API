package repository

import (
	"expense-tracker-api/domain"
	"time"
)

type categoryRepository struct{}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository() domain.CategoryRepository {
	return &categoryRepository{}
}

func (r *categoryRepository) Create(category *domain.Category) error {
	query := `INSERT INTO categories (name, created_at) 
			  VALUES ($1, $2) RETURNING id`
	err := DB.QueryRow(query, category.Name, time.Now()).Scan(&category.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *categoryRepository) GetByID(id int) (*domain.Category, error) {
	category := &domain.Category{}
	query := `SELECT id, name, created_at FROM categories WHERE id = $1`
	err := DB.QueryRow(query, id).Scan(&category.ID, &category.Name, &category.CreatedAt)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (r *categoryRepository) GetAll() ([]*domain.Category, error) {
	query := `SELECT id, name, created_at FROM categories ORDER BY name`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		category := &domain.Category{}
		err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (r *categoryRepository) Update(category *domain.Category) error {
	query := `UPDATE categories SET name = $1 WHERE id = $2`
	_, err := DB.Exec(query, category.Name, category.ID)
	return err
}

func (r *categoryRepository) Delete(id int) error {
	query := `DELETE FROM categories WHERE id = $1`
	_, err := DB.Exec(query, id)
	return err
}
