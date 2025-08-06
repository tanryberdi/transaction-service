# Transaction Service

A high-performance transaction processing service built with Go, implementing hexagonal architecture with Gin web framework and PostgreSQL database.

## Features

- **Hexagonal Architecture**: Clean separation of concerns with domain, application, and infrastructure layers
- **Concurrent Processing**: Handles 20-30+ requests per second with proper concurrent safety
- **Idempotent Transactions**: Prevents duplicate transaction processing using unique transaction IDs
- **Balance Management**: Maintains user account balances with precise decimal arithmetic
- **Docker Ready**: Complete containerization with Docker Compose for easy deployment

## Quick Start

### Prerequisites
- Docker and Docker Compose installed

### Running the Application

1. Clone the repository and navigate to the project directory

2. Start the application with Docker Compose:
```bash
docker-compose up -d
```

This command will:
- Start PostgreSQL database
- Build and run the Go application
- Create predefined users (IDs: 1, 2, 3) with initial balance of 100.00
- Expose the API on port 8080

3. The service will be available at `http://localhost:8080`

### Stopping the Application
```bash
docker-compose down
```

To remove volumes (database data):
```bash
docker-compose down -v
```

## API Endpoints

### 1. Process Transaction
**POST** `/user/{userId}/transaction`

Updates user balance based on transaction type.

**Headers:**
- `Source-Type`: `game`, `server`, or `payment`
- `Content-Type`: `application/json`

**Request Body:**
```json
{
  "state": "win",           // "win" or "lose"
  "amount": "10.15",        // string, up to 2 decimal places
  "transactionId": "uuid-123"  // unique transaction identifier
}
```

**Example Request:**
```bash
curl -X POST http://localhost:8080/user/1/transaction \
  -H "Source-Type: game" \
  -H "Content-Type: application/json" \
  -d '{
    "state": "win",
    "amount": "25.50",
    "transactionId": "tx-001"
  }'
```

**Success Response (200 OK):**
```json
{
  "message": "Transaction processed successfully"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid input data
- `404 Not Found`: User not found
- `409 Conflict`: Duplicate transaction ID

### 2. Get User Balance
**GET** `/user/{userId}/balance`

Retrieves current user balance.

**Example Request:**
```bash
curl http://localhost:8080/user/1/balance
```

**Success Response (200 OK):**
```json
{
  "userId": 1,
  "balance": "125.50"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid user ID
- `404 Not Found`: User not found

## Testing the Application

### Basic Test Scenarios

1. **Check initial balance:**
```bash
curl http://localhost:8080/user/1/balance
# Expected: {"userId":1,"balance":"100.00"}
```

2. **Process a winning transaction:**
```bash
curl -X POST http://localhost:8080/user/1/transaction \
  -H "Source-Type: game" \
  -H "Content-Type: application/json" \
  -d '{
    "state": "win",
    "amount": "25.50",
    "transactionId": "tx-win-001"
  }'
```

3. **Check updated balance:**
```bash
curl http://localhost:8080/user/1/balance
# Expected: {"userId":1,"balance":"125.50"}
```

4. **Process a losing transaction:**
```bash
curl -X POST http://localhost:8080/user/1/transaction \
  -H "Source-Type: game" \
  -H "Content-Type: application/json" \
  -d '{
    "state": "lose",
    "amount": "15.25",
    "transactionId": "tx-lose-001"
  }'
```

5. **Check final balance:**
```bash
curl http://localhost:8080/user/1/balance
# Expected: {"userId":1,"balance":"110.25"}
```

6. **Test duplicate transaction (should fail):**
```bash
curl -X POST http://localhost:8080/user/1/transaction \
  -H "Source-Type: game" \
  -H "Content-Type: application/json" \
  -d '{
    "state": "win",
    "amount": "10.00",
    "transactionId": "tx-win-001"
  }'
# Expected: 409 Conflict - Transaction already processed
```

### Load Testing

You can use tools like `ab` (Apache Bench) or `hey` to test concurrent performance:

```bash
# Install hey: go install github.com/rakyll/hey@latest

# Test with 30 concurrent requests
hey -n 1000 -c 30 -m POST \
  -H "Source-Type: game" \
  -H "Content-Type: application/json" \
  -d '{"state":"win","amount":"1.00","transactionId":"load-test-"}' \
  http://localhost:8080/user/1/transaction
```

## Architecture

The application follows **Hexagonal Architecture** (Ports and Adapters) pattern:

### Layer Structure

1. **Domain Layer** (`internal/domain/`)
    - **Entities**: Core business objects (User, Transaction)
    - **Repositories**: Interfaces for data access

2. **Application Layer** (`internal/application/`)
    - **Services**: Business logic and use cases
    - **Transaction processing with validation and balance calculations**

3. **Infrastructure Layer** (`internal/adapters/`)
    - **Database**: PostgreSQL implementation of repositories
    - **Handlers**: Gin web framework handlers and routing

### Key Design Decisions

- **Decimal Precision**: Uses `shopspring/decimal` library for accurate financial calculations
- **Idempotency**: Ensures each transaction ID is processed only once
- **Concurrent Safety**: Database transactions prevent race conditions
- **Error Handling**: Comprehensive error types with appropriate HTTP status codes
- **Connection Pooling**: Optimized database connection management

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    balance DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Transactions Table
```sql
CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    transaction_id VARCHAR(255) NOT NULL UNIQUE,
    state VARCHAR(10) NOT NULL CHECK (state IN ('win', 'lose')),
    amount DECIMAL(15,2) NOT NULL,
    source_type VARCHAR(20) NOT NULL CHECK (source_type IN ('game', 'server', 'payment')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Development

### Local Development Setup

1. **Install dependencies:**
```bash
go mod download
```

2. **Start PostgreSQL (using Docker):**
```bash
docker run --name postgres -e POSTGRES_PASSWORD=tanryberdi -e POSTGRES_DB=transaction -p 5432:5432 -d postgres:16-alpine
```

3. **Set environment variables:**
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=tanryberdi
export DB_PASSWORD=tanryberdi
export DB_NAME=transaction
export DB_SSLMODE=disable
export PORT=8080
```

4. **Run the application:**
```bash
go run main.go
```

### Project Structure
```
transaction-service/
├── main.go                          # Application entry point
├── go.mod                           # Go module definition
├── go.sum                           # Go module checksums
├── Dockerfile                       # Docker container definition
├── docker-compose.yml               # Docker Compose configuration
├── .env                            # Environment variables
├── README.md                       # This file
└── internal/
    ├── domain/
    │   ├── entities/
    │   │   └── entities.go         # Domain entities
    │   └── repositories/
    │       └── repositories.go     # Repository interfaces
    ├── application/
    │   └── services/
    │       └── transaction_service.go  # Business logic
    └── adapters/
        ├── database/
        │   ├── connection.go       # Database connection
        │   ├── migrations.go       # Database migrations
        │   ├── user_repository.go  # User repository implementation
        │   └── transaction_repository.go  # Transaction repository implementation
        └── handlers/
            └── handlers.go         # HTTP handlers
```

## Performance Considerations

- **Connection Pooling**: Database connections are pooled (max 25 open, 5 idle)
- **Indexing**: Proper database indexes for fast lookups
- **Decimal Arithmetic**: Precise financial calculations without floating-point errors
- **Concurrent Safety**: Proper transaction isolation for concurrent requests
- **Memory Efficiency**: Minimal memory footprint with efficient data structures

## Security Considerations

- **Input Validation**: Comprehensive validation of all input parameters
- **SQL Injection Prevention**: Using parameterized queries
- **Rate Limiting**: Ready for rate limiting middleware addition
- **Error Handling**: No sensitive information exposed in error messages

## Monitoring and Observability

The application includes:
- **Structured Logging**: Request/response logging via Gin middleware
- **Health Checks**: Database connectivity verification
- **Error Tracking**: Comprehensive error handling and reporting
- **Metrics Ready**: Prepared for metrics collection integration