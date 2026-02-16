package repository

import (
	"expense-tracker-api/domain"
	"time"
)

type budgetRepository struct{}

// NewBudgetRepository creates a new budget repository
func NewBudgetRepository() domain.BudgetRepository {
	return &budgetRepository{}
}

func (r *budgetRepository) Create(budget *domain.Budget) error {
	query := `INSERT INTO budgets (month, year, budget_amount, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	now := time.Now()
	err := DB.QueryRow(query, budget.Month, budget.Year, budget.BudgetAmount, now, now).Scan(&budget.ID)
	if err != nil {
		return err
	}
	budget.CreatedAt = now
	budget.UpdatedAt = now
	return nil
}

func (r *budgetRepository) GetByID(id int) (*domain.Budget, error) {
	budget := &domain.Budget{}
	query := `SELECT id, month, year, budget_amount, created_at, updated_at FROM budgets WHERE id = $1`
	err := DB.QueryRow(query, id).Scan(&budget.ID, &budget.Month, &budget.Year,
		&budget.BudgetAmount, &budget.CreatedAt, &budget.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return budget, nil
}

func (r *budgetRepository) GetAll() ([]*domain.Budget, error) {
	query := `SELECT id, month, year, budget_amount, created_at, updated_at 
			  FROM budgets ORDER BY year DESC, month DESC`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var budgets []*domain.Budget
	for rows.Next() {
		budget := &domain.Budget{}
		err := rows.Scan(&budget.ID, &budget.Month, &budget.Year,
			&budget.BudgetAmount, &budget.CreatedAt, &budget.UpdatedAt)
		if err != nil {
			return nil, err
		}
		budgets = append(budgets, budget)
	}
	return budgets, nil
}

func (r *budgetRepository) GetByMonth(month, year int) (*domain.Budget, error) {
	budget := &domain.Budget{}
	query := `SELECT id, month, year, budget_amount, created_at, updated_at 
			  FROM budgets WHERE month = $1 AND year = $2`
	err := DB.QueryRow(query, month, year).Scan(&budget.ID, &budget.Month, &budget.Year,
		&budget.BudgetAmount, &budget.CreatedAt, &budget.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return budget, nil
}

func (r *budgetRepository) Update(budget *domain.Budget) error {
	query := `UPDATE budgets SET budget_amount = $1, updated_at = $2 WHERE id = $3`
	budget.UpdatedAt = time.Now()
	_, err := DB.Exec(query, budget.BudgetAmount, budget.UpdatedAt, budget.ID)
	return err
}

func (r *budgetRepository) Delete(id int) error {
	query := `DELETE FROM budgets WHERE id = $1`
	_, err := DB.Exec(query, id)
	return err
}
