package domain

import "time"

// Budget represents a monthly budget
type Budget struct {
	ID           int       `json:"id"`
	Month        int       `json:"month"`
	Year         int       `json:"year"`
	BudgetAmount float64   `json:"budget_amount"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// BudgetStatus represents the status of a budget with spending information
type BudgetStatus struct {
	Budget      *Budget `json:"budget"`
	SpentAmount float64 `json:"spent_amount"`
	Remaining   float64 `json:"remaining"`
	Status      string  `json:"status"` // "within_budget" or "exceeded"
}

// BudgetRepository defines the interface for budget data operations
type BudgetRepository interface {
	Create(budget *Budget) error
	GetByID(id int) (*Budget, error)
	GetAll() ([]*Budget, error)
	GetByMonth(month, year int) (*Budget, error)
	Update(budget *Budget) error
	Delete(id int) error
}
