package transport

import (
	"expense-tracker-api/services"
	"expense-tracker-api/transport/handlers"
	"expense-tracker-api/transport/middleware"

	"github.com/gorilla/mux"
)

// SetupRouter sets up all routes
func SetupRouter(categoryService *services.CategoryService,
	expenseService *services.ExpenseService, budgetService *services.BudgetService) *mux.Router {

	router := mux.NewRouter()

	// Initialize handlers
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	expenseHandler := handlers.NewExpenseHandler(expenseService)
	budgetHandler := handlers.NewBudgetHandler(budgetService)

	// Apply CORS middleware
	router.Use(middleware.CORSMiddleware)

	// API routes
	api := router.PathPrefix("/api").Subrouter()

	// Category routes
	api.HandleFunc("/categories", categoryHandler.GetCategories).Methods("GET", "OPTIONS")
	api.HandleFunc("/categories", categoryHandler.CreateCategory).Methods("POST", "OPTIONS")
	api.HandleFunc("/categories/{id}", categoryHandler.UpdateCategory).Methods("PUT", "OPTIONS")
	api.HandleFunc("/categories/{id}", categoryHandler.DeleteCategory).Methods("DELETE", "OPTIONS")

	// Expense routes
	api.HandleFunc("/expenses", expenseHandler.GetExpenses).Methods("GET", "OPTIONS")
	api.HandleFunc("/expenses/{id}", expenseHandler.GetExpense).Methods("GET", "OPTIONS")
	api.HandleFunc("/expenses", expenseHandler.CreateExpense).Methods("POST", "OPTIONS")
	api.HandleFunc("/expenses/{id}", expenseHandler.UpdateExpense).Methods("PUT", "OPTIONS")
	api.HandleFunc("/expenses/{id}", expenseHandler.DeleteExpense).Methods("DELETE", "OPTIONS")

	// Budget routes
	api.HandleFunc("/budgets", budgetHandler.GetBudgets).Methods("GET", "OPTIONS")
	api.HandleFunc("/budgets/{month}/{year}", budgetHandler.GetBudgetByMonth).Methods("GET", "OPTIONS")
	api.HandleFunc("/budgets", budgetHandler.CreateOrUpdateBudget).Methods("POST", "OPTIONS")
	api.HandleFunc("/budgets/{id}", budgetHandler.DeleteBudget).Methods("DELETE", "OPTIONS")

	return router
}
