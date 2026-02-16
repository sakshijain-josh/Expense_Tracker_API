package domain

import "time"

// Expense represents an expense entry
type Expense struct {
	ID          int         `json:"id"`
	CategoryID  int         `json:"category_id"`
	Amount      float64     `json:"amount"`
	Description string      `json:"description"`
	PaymentMode PaymentMode `json:"payment_mode"`
	ExpenseDate time.Time   `json:"expense_date"`
	CreatedAt   time.Time   `json:"created_at"`
	Warning     string      `json:"warning,omitempty"`
}

// ExpenseFilter represents filters for querying expenses
type ExpenseFilter struct {
	CategoryID  *int
	PaymentMode *PaymentMode
	StartDate   *time.Time
	EndDate     *time.Time
}

// ExpenseRepository defines the interface for expense data operations
type ExpenseRepository interface {
	Create(expense *Expense) error
	GetByID(id int) (*Expense, error)
	GetAll(filter *ExpenseFilter) ([]*Expense, error)
	Update(expense *Expense) error
	Delete(id int) error
	GetTotalByMonth(month, year int) (float64, error)
}
