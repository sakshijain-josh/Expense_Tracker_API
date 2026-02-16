package repository

import (
	"expense-tracker-api/domain"
	"fmt"
	"time"
)

type expenseRepository struct{}

// NewExpenseRepository creates a new expense repository
func NewExpenseRepository() domain.ExpenseRepository {
	return &expenseRepository{}
}

func (r *expenseRepository) Create(expense *domain.Expense) error {
	query := `INSERT INTO expenses (category_id, amount, description, payment_mode, expense_date, created_at) 
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := DB.QueryRow(query, expense.CategoryID, expense.Amount, expense.Description,
		expense.PaymentMode, expense.ExpenseDate, time.Now()).Scan(&expense.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *expenseRepository) GetByID(id int) (*domain.Expense, error) {
	expense := &domain.Expense{}
	query := `SELECT id, category_id, amount, description, payment_mode, expense_date, created_at 
			  FROM expenses WHERE id = $1`
	err := DB.QueryRow(query, id).Scan(&expense.ID, &expense.CategoryID, &expense.Amount,
		&expense.Description, &expense.PaymentMode, &expense.ExpenseDate, &expense.CreatedAt)
	if err != nil {
		return nil, err
	}
	return expense, nil
}

func (r *expenseRepository) GetAll(filter *domain.ExpenseFilter) ([]*domain.Expense, error) {
	query := `SELECT id, category_id, amount, description, payment_mode, expense_date, created_at 
			  FROM expenses WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if filter.CategoryID != nil {
		query += fmt.Sprintf(` AND category_id = $%d`, argIndex)
		args = append(args, *filter.CategoryID)
		argIndex++
	}

	if filter.PaymentMode != nil {
		query += fmt.Sprintf(` AND payment_mode = $%d`, argIndex)
		args = append(args, string(*filter.PaymentMode))
		argIndex++
	}

	if filter.StartDate != nil {
		query += fmt.Sprintf(` AND expense_date >= $%d`, argIndex)
		args = append(args, *filter.StartDate)
		argIndex++
	}

	if filter.EndDate != nil {
		query += fmt.Sprintf(` AND expense_date <= $%d`, argIndex)
		args = append(args, *filter.EndDate)
		argIndex++
	}

	query += ` ORDER BY expense_date DESC, created_at DESC`

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []*domain.Expense
	for rows.Next() {
		expense := &domain.Expense{}
		err := rows.Scan(&expense.ID, &expense.CategoryID, &expense.Amount,
			&expense.Description, &expense.PaymentMode, &expense.ExpenseDate, &expense.CreatedAt)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}
	return expenses, nil
}

func (r *expenseRepository) Update(expense *domain.Expense) error {
	query := `UPDATE expenses SET category_id = $1, amount = $2, description = $3, 
			  payment_mode = $4, expense_date = $5 WHERE id = $6`
	_, err := DB.Exec(query, expense.CategoryID, expense.Amount, expense.Description,
		expense.PaymentMode, expense.ExpenseDate, expense.ID)
	return err
}

func (r *expenseRepository) Delete(id int) error {
	query := `DELETE FROM expenses WHERE id = $1`
	_, err := DB.Exec(query, id)
	return err
}

func (r *expenseRepository) GetTotalByMonth(month, year int) (float64, error) {
	var total float64
	query := `SELECT COALESCE(SUM(amount), 0) FROM expenses 
			  WHERE EXTRACT(MONTH FROM expense_date) = $1 AND EXTRACT(YEAR FROM expense_date) = $2`
	err := DB.QueryRow(query, month, year).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}
