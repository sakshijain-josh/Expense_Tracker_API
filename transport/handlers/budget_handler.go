package handlers

import (
	"encoding/json"
	"expense-tracker-api/domain"
	"expense-tracker-api/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type BudgetHandler struct {
	budgetService *services.BudgetService
}

// NewBudgetHandler creates a new budget handler
func NewBudgetHandler(budgetService *services.BudgetService) *BudgetHandler {
	return &BudgetHandler{budgetService: budgetService}
}

type CreateBudgetRequest struct {
	Month        int     `json:"month"`
	Year         int     `json:"year"`
	BudgetAmount float64 `json:"budget_amount"`
}

// GetBudgets handles getting all budgets
func (h *BudgetHandler) GetBudgets(w http.ResponseWriter, r *http.Request) {
	budgets, err := h.budgetService.GetBudgets()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(budgets)
}

// GetBudgetByMonth handles getting a budget for a specific month with status
func (h *BudgetHandler) GetBudgetByMonth(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	month, err := strconv.Atoi(vars["month"])
	if err != nil {
		http.Error(w, "Invalid month", http.StatusBadRequest)
		return
	}

	year, err := strconv.Atoi(vars["year"])
	if err != nil {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	budgetStatus, err := h.budgetService.GetBudgetByMonth(month, year)
	if err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(budgetStatus)
}

// CreateOrUpdateBudget handles creating or updating a budget
func (h *BudgetHandler) CreateOrUpdateBudget(w http.ResponseWriter, r *http.Request) {
	var req CreateBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	budget, err := h.budgetService.CreateOrUpdateBudget(req.Month, req.Year, req.BudgetAmount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(budget)
}

// DeleteBudget handles deleting a budget
func (h *BudgetHandler) DeleteBudget(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	budgetID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid budget ID", http.StatusBadRequest)
		return
	}

	err = h.budgetService.DeleteBudget(budgetID)
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
