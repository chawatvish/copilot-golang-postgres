# Gin Simple REST API with PostgreSQL

A simple REST API built with Go using the Gin web framework and GORM for PostgreSQL database integration.

## Recent Updates

### Version 2.0 - User Contact Information (August 2025)

Added comprehensive contact information support to the User model:

- **Phone Number Field**: Required in API requests, stored as nullable text in database
- **Address Field**: Optional in both API and database
- **Database Migration**: Seamless upgrade of existing tables with backward compatibility
- **Pointer Types**: Uses Go pointer types (`*string`) to handle nullable database fields
- **API Validation**: Phone number required for new users, address optional
- **Test Coverage**: All existing tests pass, new functionality fully tested

#### Migration Details

- Existing users automatically get `NULL` values for new fields
- No data loss during migration
- Database and in-memory modes both supported
- GORM handles the migration automatically with proper NULL constraints

#### API Changes

- **CREATE/UPDATE requests**: Now require `phone` field, `address` is optional
- **Response format**: Includes phone and address fields in all user responses
- **Backward compatibility**: Existing API clients need to provide phone field

For detailed technical documentation, see [User Model Update Documentation](docs/user-model-update.md).

## Features

- Clean Architecture pattern with layered design
- RESTful endpoints for user management (CRUD operations)
- User model with contact information (phone number and address)
- PostgreSQL database integration with GORM
- In-memory fallback for testing and development
- Database migrations with backward compatibility
- Comprehensive test suite
- Environment-based configuration
- Database migrations and seeding
- JSON structured responses
- Error handling with proper HTTP status codes
- Pointer type handling for nullable database fields

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── database/
│   │   └── database.go          # Database connection and migrations
│   ├── handlers/
│   │   ├── health_handler.go    # Health check handlers
│   │   └── user_handler.go      # User-related handlers
│   ├── models/
│   │   └── user.go              # Data models with GORM tags
│   ├── repository/
│   │   ├── interfaces.go        # Repository interfaces
│   │   ├── memory_user_repository.go  # In-memory implementation
│   │   └── user_repository.go   # PostgreSQL implementation with GORM
│   ├── router/
│   │   └── router.go            # Route definitions and middleware
│   └── services/
│       └── user_service.go      # Business logic layer
├── pkg/
│   └── response/
│       └── response.go          # Standardized API response structure
├── tests/
│   └── api_test.go              # Comprehensive test suite
├── .env                         # Environment variables (create from .env.example)
├── go.mod                       # Go module dependencies
├── go.sum                       # Go module checksums
├── Makefile                     # Build and run commands
├── test.sh                      # Test runner script
└── README.md                    # This file
```

## Prerequisites

- Go 1.19 or higher
- PostgreSQL 12 or higher
- Git

## Installation

1. Clone the repository:

```bash
git clone <repository-url>
cd gin-simple-app
```

2. Install dependencies:

```bash
go mod tidy
```

3. Set up environment variables:

```bash
cp .env.example .env
```

Edit the `.env` file with your database credentials:

```env
# Server Configuration
PORT=8080
GIN_MODE=release

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=gin_simple_db
DB_SSLMODE=disable
DB_TIMEZONE=UTC

# Application Configuration
USE_DATABASE=true
```

## Database Setup

1. Create a PostgreSQL database:

```sql
CREATE DATABASE gin_simple_db;
```

2. The application will automatically:
   - Connect to the database on startup
   - Run migrations to create the `users` table
   - Seed initial test data
   - Fall back to in-memory storage if database connection fails

## Running the Application

### Using Make (recommended)

```bash
# Run the application
make run

# Run tests
make test

# Build the application
make build

# Clean build artifacts
make clean
```

### Manual Commands

```bash
# Run the application
go run cmd/server/main.go

# Run tests
go test ./tests/... -v

# Build the application
go build -o bin/server cmd/server/main.go
```

## API Endpoints

### Health Check

```http
GET /health
```

Returns application health status.

### Root

```http
GET /
```

Returns welcome message and API version.

### Users

#### Get All Users

```http
GET /api/v1/users
```

#### Get User by ID

```http
GET /api/v1/users/{id}
```

#### Create User

```http
POST /api/v1/users
Content-Type: application/json

{
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+1-555-0123",
    "address": "123 Main Street"
}
```

Note: `phone` is required, `address` is optional.

#### Update User

```http
PUT /api/v1/users/{id}
Content-Type: application/json

{
    "name": "John Updated",
    "email": "john.updated@example.com",
    "phone": "+1-555-0124",
    "address": "456 Oak Avenue"
}
```

#### Delete User

```http
DELETE /api/v1/users/{id}
```

## Response Format

All API responses follow this structure:

```json
{
  "success": true,
  "message": "Operation successful",
  "data": {
    /* response data */
  },
  "count": 1,
  "timestamp": "2024-01-01T12:00:00Z"
}
```

Error responses:

```json
{
  "success": false,
  "error": "Error description",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

## Testing

The application includes a comprehensive test suite that covers:

- Health check endpoint
- All user CRUD operations
- Error handling scenarios
- Complete user lifecycle testing
- Input validation
- Duplicate email constraint validation

Run tests:

```bash
# Using the test script
./test.sh

# Using Go directly
go test ./tests/... -v

# Using Make
make test
```

## Configuration

The application supports two modes:

1. **Database Mode** (default): Uses PostgreSQL with GORM

   - Set `USE_DATABASE=true` in `.env`
   - Requires valid database connection parameters

2. **In-Memory Mode**: Uses in-memory storage for testing
   - Set `USE_DATABASE=false` in `.env`
   - No database required
   - Automatic fallback if database connection fails

## Database Schema

### Users Table

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone TEXT,
    address TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
```

#### User Fields

- `name`: Required string, user's full name
- `email`: Required string, unique email address
- `phone`: Optional string in database, required in API (nullable for existing records)
- `address`: Optional string, user's physical address
- Standard GORM timestamps (created_at, updated_at, deleted_at)

## Development

### Adding New Endpoints

1. Define the handler in `internal/handlers/`
2. Add business logic in `internal/services/`
3. Update repository interface in `internal/repository/interfaces.go`
4. Implement repository methods in both GORM and in-memory implementations
5. Register routes in `internal/router/router.go`
6. Add tests in `tests/`

### Database Migrations

GORM handles migrations automatically using `AutoMigrate`. For custom migrations:

1. Add migration logic in `internal/database/database.go`
2. Update the `InitDatabase` function
3. Run the application to apply migrations

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

This project is licensed under the MIT License.

## Architecture Decisions

### Clean Architecture

The application follows clean architecture principles with clear separation of concerns:

- **Handlers**: HTTP request/response handling
- **Services**: Business logic and validation
- **Repository**: Data access abstraction
- **Models**: Data structures and validation rules

### Dependency Injection

All components are injected as dependencies, making the code:

- Testable (easy to mock dependencies)
- Flexible (easy to swap implementations)
- Maintainable (clear dependencies)

### Interface-Based Design

Repository pattern with interfaces allows:

- Multiple implementations (GORM + in-memory)
- Easy testing with mock implementations
- Future database technology changes

### Error Handling

Consistent error handling with:

- Proper HTTP status codes
- Structured error responses
- Database error translation
- Validation error mapping
