# Expense Tracker REST API

A personal expense tracker REST API built with Go, featuring expense management, categories, payment modes (UPI/Cash), and monthly budget tracking.

## Features

- **Expense Management**: Full CRUD operations for expenses
- **Categories**: Expense categories for organizing expenses
- **Payment Modes**: Track expenses by payment mode (UPI or Cash)
- **Monthly Budgets**: Set and track monthly budgets with status (within budget/exceeded)
- **Filtering**: Filter expenses by date range, category, and payment mode

## Prerequisites

- Go 1.24 or higher
- PostgreSQL 12 or higher
- Postman (for API testing)

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd Expense_Tracker_APIs
```

2. Install dependencies:
```bash
go mod download
```

3. Set up PostgreSQL database:
```bash
createdb expense_tracker
```

4. Create a `.env` file from `.env.example`:
```bash
cp .env.example .env
```

5. Update `.env` with your database credentials:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=expense_tracker
PORT=8080
```

6. Run the application:
```bash
go run main.go
```

The server will start on port 8080 (or the port specified in `.env`).

## Database Schema

- **categories**: id, name, created_at
- **expenses**: id, category_id, amount, description, payment_mode, expense_date, created_at
- **budgets**: id, month, year, budget_amount, created_at, updated_at

