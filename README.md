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

## API Endpoints

All endpoints are public and do not require authentication.

### Categories

#### Get All Categories
- **GET** `/api/categories`
- **Response**: Array of category objects
- **Example Response**:
```json
[
  {
    "id": 1,
    "name": "Food",
    "created_at": "2024-01-15T10:00:00Z"
  }
]
```

#### Create Category
- **POST** `/api/categories`
- **Request Body**:
```json
{
  "name": "Food"
}
```
- **Response**: Created category object

#### Update Category
- **PUT** `/api/categories/{id}`
- **Request Body**:
```json
{
  "name": "Groceries"
}
```
- **Response**: Updated category object

#### Delete Category
- **DELETE** `/api/categories/{id}`
- **Response**: 204 No Content

### Expenses

#### Get Expenses
- **GET** `/api/expenses`
- **Query Parameters** (optional):
  - `category_id`: Filter by category ID
  - `payment_mode`: Filter by payment mode (UPI or Cash)
  - `start_date`: Filter from date (YYYY-MM-DD)
  - `end_date`: Filter to date (YYYY-MM-DD)
- **Example**: `/api/expenses?category_id=1&payment_mode=UPI&start_date=2024-01-01&end_date=2024-01-31`
- **Response**: Array of expense objects

#### Get Single Expense
- **GET** `/api/expenses/{id}`
- **Response**: Expense object

#### Create Expense
- **POST** `/api/expenses`
- **Request Body**:
```json
{
  "category_id": 1,
  "amount": 500.50,
  "description": "Lunch at restaurant",
  "payment_mode": "UPI",
  "expense_date": "2024-01-15"
}
```
- **Note**: `expense_date` is optional (defaults to current date)
- **Response**: Created expense object

#### Update Expense
- **PUT** `/api/expenses/{id}`
- **Request Body** (all fields optional):
```json
{
  "category_id": 2,
  "amount": 600.00,
  "description": "Updated description",
  "payment_mode": "Cash",
  "expense_date": "2024-01-16"
}
```
- **Response**: Updated expense object

#### Delete Expense
- **DELETE** `/api/expenses/{id}`
- **Response**: 204 No Content

### Budgets

#### Get All Budgets
- **GET** `/api/budgets`
- **Response**: Array of budget objects
- **Example Response**:
```json
[
  {
    "id": 1,
    "month": 1,
    "year": 2024,
    "budget_amount": 5000.00,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
]
```

#### Get Budget by Month (with Status)
- **GET** `/api/budgets/{month}/{year}`
- **Example**: `/api/budgets/1/2024`
- **Response**:
```json
{
  "budget": {
    "id": 1,
    "month": 1,
    "year": 2024,
    "budget_amount": 5000.00
  },
  "spent_amount": 3500.00,
  "remaining": 1500.00,
  "status": "within_budget"
}
```
- **Status**: `"within_budget"` or `"exceeded"`

#### Create or Update Budget
- **POST** `/api/budgets`
- **Request Body**:
```json
{
  "month": 1,
  "year": 2024,
  "budget_amount": 5000.00
}
```
- **Note**: If a budget already exists for the month/year, it will be updated.
- **Response**: Created or updated budget object

#### Delete Budget
- **DELETE** `/api/budgets/{id}`
- **Response**: 204 No Content

## Postman Setup

1. **Import Collection**: Create a new collection in Postman for testing

2. **Set Environment Variables**:
   - Create a new environment
   - Add variable `base_url` = `http://localhost:8080`

3. **Example Request Flow**:
   ```
   1. POST /api/categories (create category)
   2. POST /api/expenses (create expense)
   3. POST /api/budgets (set budget)
   4. GET /api/budgets/1/2024 (check budget status)
   5. GET /api/expenses?category_id=1 (filter expenses)
   ```

## Testing

Run unit tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## Project Structure

```
Expense_Tracker_APIs/
├── domain/              # Domain models and interfaces
│   ├── expense.go
│   ├── category.go
│   ├── budget.go
│   ├── payment_mode.go
│   └── errors.go
├── repository/          # Database access layer
│   ├── db.go
│   ├── expense_repository.go
│   ├── category_repository.go
│   └── budget_repository.go
├── services/           # Business logic layer
│   ├── expense_service.go
│   ├── category_service.go
│   └── budget_service.go
├── transport/          # HTTP handlers and routing
│   ├── handlers/
│   │   ├── expense_handler.go
│   │   ├── category_handler.go
│   │   └── budget_handler.go
│   ├── middleware/
│   │   └── cors_middleware.go
│   └── router.go
├── mock/               # Mock data generators
│   └── mock_data.go
├── main.go             # Application entry point
├── go.mod
├── .env.example
└── README.md
```

## Database Schema

- **categories**: id, name, created_at
- **expenses**: id, category_id, amount, description, payment_mode, expense_date, created_at
- **budgets**: id, month, year, budget_amount, created_at, updated_at

## Error Handling

The API returns standard HTTP status codes:
- `200 OK`: Successful request
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid input
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

Error responses include a message in the response body.

## Security Notes

- All endpoints are public (no authentication required)
- Use HTTPS in production
- Implement rate limiting for production use
- Consider adding authentication if deploying to a shared environment

## License

This project is open source and available under the MIT License.
