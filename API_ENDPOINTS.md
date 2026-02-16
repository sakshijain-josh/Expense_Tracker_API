# API Endpoints Guide 

## Base URL
```
http://localhost:8080/api
```

## Important Notes
- All endpoints require `Content-Type: application/json` header for POST/PUT requests
- All endpoints are public (no authentication required)
- Date format: `YYYY-MM-DD` (e.g., "2024-01-15")
- Payment modes: `"UPI"` or `"Cash"` (case-sensitive)

---

## Categories

### 1. Get All Categories
**GET** `/api/categories`

---
### 2. Create Category
**POST** `/api/categories`

**Common Errors:**
- `400 Bad Request` - "Invalid request body" - Missing Content-Type header or invalid JSON
- `400 Bad Request` - "invalid input" - Empty name field

---

### 3. Update Category
**PUT** `/api/categories/{id}`

---

### 4. Delete Category
**DELETE** `/api/categories/{id}`

---

## Expenses

### 1. Get All Expenses
**GET** `/api/expenses`

**Query Parameters (all optional):**
- `category_id` - Filter by category ID
- `payment_mode` - Filter by payment mode (UPI or Cash)
- `start_date` - Filter from date (YYYY-MM-DD)
- `end_date` - Filter to date (YYYY-MM-DD)

**Sample Requests:**
```bash
# Get all expenses
curl http://localhost:8080/api/expenses

# Filter by category
curl "http://localhost:8080/api/expenses?category_id=1"

# Filter by payment mode
curl "http://localhost:8080/api/expenses?payment_mode=UPI"

# Filter by date range
curl "http://localhost:8080/api/expenses?start_date=2024-01-01&end_date=2024-01-31"

# Combined filters
curl "http://localhost:8080/api/expenses?category_id=1&payment_mode=UPI&start_date=2024-01-01&end_date=2024-01-31"
```

**Response:**
```json
[
  {
    "id": 1,
    "category_id": 1,
    "amount": 500.50,
    "description": "Lunch at restaurant",
    "payment_mode": "UPI",
    "expense_date": "2024-01-15T00:00:00Z",
    "created_at": "2024-01-15T10:00:00Z"
  }
]
```

---

### 2. Get Single Expense
**GET** `/api/expenses/{id}`

**Sample Request:**
```bash
curl http://localhost:8080/api/expenses/1
```

---

### 3. Create Expense
**POST** `/api/expenses`

**Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "category_id": 1,
  "amount": 500.50,
  "description": "Lunch at restaurant",
  "payment_mode": "UPI",
  "expense_date": "2024-01-15"
}
```

**Note:** `expense_date` is optional. If not provided, it defaults to current date.

**Sample Request (cURL):**
```bash
curl -X POST http://localhost:8080/api/expenses \
  -H "Content-Type: application/json" \
  -d '{
    "category_id": 1,
    "amount": 500.50,
    "description": "Lunch at restaurant",
    "payment_mode": "UPI",
    "expense_date": "2024-01-15"
  }'
```

**Minimal Request (without optional fields):**
```json
{
  "category_id": 1,
  "amount": 500.50,
  "payment_mode": "UPI"
}
```

**Response (201 Created):**
```json
{
  "id": 1,
  "category_id": 1,
  "amount": 500.50,
  "description": "Lunch at restaurant",
  "payment_mode": "UPI",
  "expense_date": "2024-01-15T00:00:00Z",
  "created_at": "2024-01-15T10:00:00Z"
}
```

**Common Errors:**
- `400 Bad Request` - "Invalid request body" - Missing Content-Type header or invalid JSON
- `400 Bad Request` - "invalid payment mode" - Payment mode must be "UPI" or "Cash"
- `400 Bad Request` - "invalid category" - Category ID does not exist

---

### 4. Update Expense
**PUT** `/api/expenses/{id}`

**Headers:**
```
Content-Type: application/json
```

**Request Body (all fields optional):**
```json
{
  "category_id": 2,
  "amount": 600.00,
  "description": "Updated description",
  "payment_mode": "Cash",
  "expense_date": "2024-01-16"
}
```

**Sample Request (cURL):**
```bash
curl -X PUT http://localhost:8080/api/expenses/1 \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 600.00,
    "description": "Updated description"
  }'
```

**Response (200 OK):**
```json
{
  "id": 1,
  "category_id": 2,
  "amount": 600.00,
  "description": "Updated description",
  "payment_mode": "Cash",
  "expense_date": "2024-01-16T00:00:00Z",
  "created_at": "2024-01-15T10:00:00Z"
}
```

---

### 5. Delete Expense
**DELETE** `/api/expenses/{id}`

**Sample Request:**
```bash
curl -X DELETE http://localhost:8080/api/expenses/1
```

**Response:** `204 No Content`

---

## Budgets

### 1. Get All Budgets
**GET** `/api/budgets`

**Sample Request:**
```bash
curl http://localhost:8080/api/budgets
```

**Response:**
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

---

### 2. Get Budget by Month (with Status)
**GET** `/api/budgets/{month}/{year}`

**Sample Request:**
```bash
curl http://localhost:8080/api/budgets/1/2024
```

**Response (200 OK):**
```json
{
  "budget": {
    "id": 1,
    "month": 1,
    "year": 2024,
    "budget_amount": 5000.00,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  },
  "spent_amount": 3500.00,
  "remaining": 1500.00,
  "status": "within_budget"
}
```

**Status values:**
- `"within_budget"` - Spending is within the budget limit
- `"exceeded"` - Spending has exceeded the budget

---

### 3. Create or Update Budget
**POST** `/api/budgets`

**Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "month": 1,
  "year": 2024,
  "budget_amount": 5000.00
}
```

**Sample Request (cURL):**
```bash
curl -X POST http://localhost:8080/api/budgets \
  -H "Content-Type: application/json" \
  -d '{
    "month": 1,
    "year": 2024,
    "budget_amount": 5000.00
  }'
```

**Note:** If a budget already exists for the month/year, it will be updated.

**Response (201 Created):**
```json
{
  "id": 1,
  "month": 1,
  "year": 2024,
  "budget_amount": 5000.00,
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:00:00Z"
}
```

**Common Errors:**
- `400 Bad Request` - "Invalid request body" - Missing Content-Type header or invalid JSON
- `400 Bad Request` - "invalid input" - Invalid month (must be 1-12) or negative budget amount

---

### 4. Delete Budget
**DELETE** `/api/budgets/{id}`

**Sample Request:**
```bash
curl -X DELETE http://localhost:8080/api/budgets/1
```

**Response:** `204 No Content`

---

## Complete Example Workflow

### Step 1: Create a Category
```bash
curl -X POST http://localhost:8080/api/categories \
  -H "Content-Type: application/json" \
  -d '{"name": "Food"}'
```
**Response:** `{"id": 1, "name": "Food", ...}`

### Step 2: Create an Expense
```bash
curl -X POST http://localhost:8080/api/expenses \
  -H "Content-Type: application/json" \
  -d '{
    "category_id": 1,
    "amount": 500.50,
    "description": "Lunch",
    "payment_mode": "UPI",
    "expense_date": "2024-01-15"
  }'
```
**Response:** `{"id": 1, "category_id": 1, "amount": 500.50, ...}`

### Step 3: Create a Budget
```bash
curl -X POST http://localhost:8080/api/budgets \
  -H "Content-Type: application/json" \
  -d '{
    "month": 1,
    "year": 2024,
    "budget_amount": 5000.00
  }'
```
**Response:** `{"id": 1, "month": 1, "year": 2024, "budget_amount": 5000.00, ...}`

### Step 4: Check Budget Status
```bash
curl http://localhost:8080/api/budgets/1/2024
```
**Response:** Shows budget with spent amount and status

---

## Postman Collection Setup

### Environment Variables
Create a new environment with:
- `base_url` = `http://localhost:8080`

### Common Headers
For POST/PUT requests, add header:
- Key: `Content-Type`
- Value: `application/json`

### Sample Postman Requests

1. **Create Category**
   - Method: POST
   - URL: `{{base_url}}/api/categories`
   - Headers: `Content-Type: application/json`
   - Body (raw JSON):
   ```json
   {
     "name": "Food"
   }
   ```

2. **Create Expense**
   - Method: POST
   - URL: `{{base_url}}/api/expenses`
   - Headers: `Content-Type: application/json`
   - Body (raw JSON):
   ```json
   {
     "category_id": 1,
     "amount": 500.50,
     "description": "Lunch at restaurant",
     "payment_mode": "UPI",
     "expense_date": "2024-01-15"
   }
   ```

3. **Create Budget**
   - Method: POST
   - URL: `{{base_url}}/api/budgets`
   - Headers: `Content-Type: application/json`
   - Body (raw JSON):
   ```json
   {
     "month": 1,
     "year": 2024,
     "budget_amount": 5000.00
   }
   ```

---
