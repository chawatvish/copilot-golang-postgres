# User Model Update: Phone and Address Fields

## Overview

This document describes the implementation of phone number and address fields to the User model in the Go REST API application.

## Changes Made

### 1. Database Schema Changes

#### New Columns Added

- `phone`: `TEXT` type, nullable, with default value `NULL`
- `address`: `TEXT` type, nullable

#### Migration Strategy

The migration was designed to handle existing data gracefully:

- Used GORM's auto-migration feature
- Added columns as nullable to avoid conflicts with existing records
- Existing users will have `NULL` values for both new fields

### 2. Model Structure Updates

#### User Struct

```go
type User struct {
    ID        uint           `json:"id" gorm:"primarykey"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
    Name      string         `json:"name" gorm:"not null" binding:"required"`
    Email     string         `json:"email" gorm:"uniqueIndex;not null" binding:"required,email"`
    Phone     *string        `json:"phone" gorm:"type:text;default:null" binding:"required"`
    Address   *string        `json:"address,omitempty" gorm:"type:text"`
}
```

#### Key Design Decisions

- **Phone Field**: Uses `*string` (pointer) to handle nullable database storage while maintaining API requirement
- **Address Field**: Uses `*string` (pointer) and is optional in both API and database
- **GORM Tags**: Explicit `type:text;default:null` to ensure proper migration behavior

### 3. API Contract Updates

#### Request Structures

```go
type CreateUserRequest struct {
    Name    string  `json:"name" binding:"required"`
    Email   string  `json:"email" binding:"required,email"`
    Phone   string  `json:"phone" binding:"required"`
    Address *string `json:"address,omitempty"`
}

type UpdateUserRequest struct {
    Name    string  `json:"name" binding:"required"`
    Email   string  `json:"email" binding:"required,email"`
    Phone   string  `json:"phone" binding:"required"`
    Address *string `json:"address,omitempty"`
}
```

#### API Behavior

- **Phone**: Required field in API requests (validation enforced)
- **Address**: Optional field (can be omitted or set to null)
- **Response**: Both fields included in JSON responses (phone as string, address as string or null)

### 4. Service Layer Changes

#### Pointer Conversion

Service methods now convert string requests to pointer fields:

```go
user := &models.User{
    Name:    req.Name,
    Email:   req.Email,
    Phone:   &req.Phone,  // Convert string to *string
    Address: req.Address, // Already *string
}
```

### 5. Repository Layer Updates

#### GORM Repository

- No changes needed - GORM handles pointer types automatically
- Migrations applied through `db.AutoMigrate(&models.User{})`

#### In-Memory Repository

- Updated test data with pointer variables
- Added sample phone numbers and addresses for testing

### 6. Database Migration Process

#### Successful Migration Log

```
ALTER TABLE "users" ADD "phone" text DEFAULT null
ALTER TABLE "users" ADD "address" text
```

#### Migration Safety

- No data loss occurred
- Existing records preserved with NULL values for new fields
- Backward compatibility maintained

## Testing

### Test Coverage

- ✅ All existing tests continue to pass
- ✅ In-memory mode tests pass
- ✅ Database mode tests pass
- ✅ API endpoints work with new fields

### Test Data

Sample data includes:

- Users with phone numbers: "+1-555-0101", "+1-555-0102", "+1-555-0103"
- Users with addresses: "123 Main St", "456 Oak Ave", "789 Pine Rd"

## API Examples

### Create User with Phone and Address

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+1-555-0123",
    "address": "123 Main Street"
  }'
```

### Create User with Phone Only

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe",
    "email": "jane@example.com",
    "phone": "+1-555-0124"
  }'
```

### Response Format

```json
{
  "id": 1,
  "created_at": "2025-08-14T22:00:00Z",
  "updated_at": "2025-08-14T22:00:00Z",
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "+1-555-0123",
  "address": "123 Main Street"
}
```

## Troubleshooting

### Common Issues Resolved

#### Migration Error: "column contains null values"

**Problem**: Initial attempt to add NOT NULL constraint failed due to existing data.
**Solution**: Used nullable field with pointer types and explicit GORM tags.

#### API Validation vs Database Storage

**Problem**: Need phone as required in API but nullable in database.
**Solution**: Used `*string` for database storage with `binding:"required"` for API validation.

## Future Considerations

### Potential Enhancements

- Phone number format validation
- Address structure with separate fields (street, city, state, zip)
- International phone number support
- Address geocoding integration

### Database Optimization

- Consider adding indexes on phone field if search functionality is needed
- Address normalization for better querying

## Dependencies

No new external dependencies were added. The implementation uses:

- Existing GORM ORM functionality
- Gin framework validation features
- Go's built-in pointer types
