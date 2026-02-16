package handlers

import (
	"encoding/json"
	"expense-tracker-api/domain"
	"expense-tracker-api/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type ExpenseHandler struct {
	expenseService *services.ExpenseService
}

// NewExpenseHandler creates a new expense handler
func NewExpenseHandler(expenseService *services.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{expenseService: expenseService}
}

type CreateExpenseRequest struct {
	CategoryID  int     `json:"category_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	PaymentMode string  `json:"payment_mode"`
	ExpenseDate string  `json:"expense_date"`
}

type UpdateExpenseRequest struct {
	CategoryID  *int     `json:"category_id"`
	Amount      *float64 `json:"amount"`
	Description *string  `json:"description"`
	PaymentMode *string  `json:"payment_mode"`
	ExpenseDate *string  `json:"expense_date"`
}

// GetExpenses handles getting expenses with optional filters
func (h *ExpenseHandler) GetExpenses(w http.ResponseWriter, r *http.Request) {
	filter := &domain.ExpenseFilter{}

	// Parse query parameters
	if categoryIDStr := r.URL.Query().Get("category_id"); categoryIDStr != "" {
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err == nil {
			filter.CategoryID = &categoryID
		}
	}

	if paymentModeStr := r.URL.Query().Get("payment_mode"); paymentModeStr != "" {
		pm := domain.PaymentMode(paymentModeStr)
		if pm.IsValid() {
			filter.PaymentMode = &pm
		}
	}

	if startDateStr := r.URL.Query().Get("start_date"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			filter.StartDate = &startDate
		}
	}

	if endDateStr := r.URL.Query().Get("end_date"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			filter.EndDate = &endDate
		}
	}

	expenses, err := h.expenseService.GetExpenses(filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

// GetExpense handles getting a single expense
func (h *ExpenseHandler) GetExpense(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	expenseID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid expense ID", http.StatusBadRequest)
		return
	}

	expense, err := h.expenseService.GetExpenseByID(expenseID)
	if err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expense)
}

// CreateExpense handles creating a new expense
func (h *ExpenseHandler) CreateExpense(w http.ResponseWriter, r *http.Request) {
	var req CreateExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	expense := &domain.Expense{
		CategoryID:  req.CategoryID,
		Amount:      req.Amount,
		Description: req.Description,
		PaymentMode: domain.PaymentMode(req.PaymentMode),
	}

	if req.ExpenseDate != "" {
		expenseDate, err := time.Parse("2006-01-02", req.ExpenseDate)
		if err == nil {
			expense.ExpenseDate = expenseDate
		}
	}

	createdExpense, err := h.expenseService.CreateExpense(expense)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdExpense)
}

// UpdateExpense handles updating an expense
func (h *ExpenseHandler) UpdateExpense(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	expenseID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid expense ID", http.StatusBadRequest)
		return
	}

	var req UpdateExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	expense := &domain.Expense{ID: expenseID}
	if req.CategoryID != nil {
		expense.CategoryID = *req.CategoryID
	}
	if req.Amount != nil {
		expense.Amount = *req.Amount
	}
	if req.Description != nil {
		expense.Description = *req.Description
	}
	if req.PaymentMode != nil {
		expense.PaymentMode = domain.PaymentMode(*req.PaymentMode)
	}
	if req.ExpenseDate != nil {
		expenseDate, err := time.Parse("2006-01-02", *req.ExpenseDate)
		if err == nil {
			expense.ExpenseDate = expenseDate
		}
	}

	updatedExpense, err := h.expenseService.UpdateExpense(expense)
	if err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedExpense)
}

// DeleteExpense handles deleting an expense
func (h *ExpenseHandler) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	expenseID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid expense ID", http.StatusBadRequest)
		return
	}

	err = h.expenseService.DeleteExpense(expenseID)
	if err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
