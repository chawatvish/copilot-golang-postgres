# API Documentation

## Base URL
```
http://localhost:8080
```

## Authentication
Currently no authentication is required for any endpoints.

## Response Format

### Success Response
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { /* response data */ },
  "count": 1,
  "timestamp": "2025-08-14T22:00:00Z"
}
```

### Error Response
```json
{
  "success": false,
  "error": "Error description",
  "timestamp": "2025-08-14T22:00:00Z"
}
```

## Endpoints

### Health Check
**GET** `/health`

Returns the health status of the API.

**Response:**
```json
{
  "success": true,
  "message": "Service is healthy",
  "data": {
    "status": "healthy",
    "database": "connected",
    "timestamp": "2025-08-14T22:00:00Z"
  }
}
```

### Root Endpoint
**GET** `/`

Returns welcome message and API information.

## User Management

### Get All Users
**GET** `/api/v1/users`

Returns a list of all users.

**Response:**
```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "phone": "+1-555-0123",
      "address": "123 Main Street",
      "created_at": "2025-08-14T22:00:00Z",
      "updated_at": "2025-08-14T22:00:00Z"
    }
  ],
  "count": 1,
  "timestamp": "2025-08-14T22:00:00Z"
}
```

### Get User by ID
**GET** `/api/v1/users/{id}`

Returns a specific user by ID.

**Parameters:**
- `id` (path, required): User ID

**Response:**
```json
{
  "success": true,
  "message": "User retrieved successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+1-555-0123",
    "address": "123 Main Street",
    "created_at": "2025-08-14T22:00:00Z",
    "updated_at": "2025-08-14T22:00:00Z"
  },
  "timestamp": "2025-08-14T22:00:00Z"
}
```

**Error Response (404):**
```json
{
  "success": false,
  "error": "User not found",
  "timestamp": "2025-08-14T22:00:00Z"
}
```

### Create User
**POST** `/api/v1/users`

Creates a new user.

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "+1-555-0123",
  "address": "123 Main Street"
}
```

**Field Requirements:**
- `name`: Required string, user's full name
- `email`: Required string, must be valid email format and unique
- `phone`: Required string, phone number
- `address`: Optional string, physical address (can be omitted)

**Response (201):**
```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "id": 4,
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+1-555-0123",
    "address": "123 Main Street",
    "created_at": "2025-08-14T22:00:00Z",
    "updated_at": "2025-08-14T22:00:00Z"
  },
  "timestamp": "2025-08-14T22:00:00Z"
}
```

**Error Response (400) - Validation Error:**
```json
{
  "success": false,
  "error": "Key: 'CreateUserRequest.Phone' Error:Tag: 'required' Tag: 'required'",
  "timestamp": "2025-08-14T22:00:00Z"
}
```

**Error Response (400) - Duplicate Email:**
```json
{
  "success": false,
  "error": "Email already exists",
  "timestamp": "2025-08-14T22:00:00Z"
}
```

### Update User
**PUT** `/api/v1/users/{id}`

Updates an existing user.

**Parameters:**
- `id` (path, required): User ID

**Request Body:**
```json
{
  "name": "John Updated",
  "email": "john.updated@example.com",
  "phone": "+1-555-0124",
  "address": "456 Oak Avenue"
}
```

**Field Requirements:**
- Same as Create User endpoint
- All fields are required in the request body

**Response (200):**
```json
{
  "success": true,
  "message": "User updated successfully",
  "data": {
    "id": 1,
    "name": "John Updated",
    "email": "john.updated@example.com",
    "phone": "+1-555-0124",
    "address": "456 Oak Avenue",
    "created_at": "2025-08-14T22:00:00Z",
    "updated_at": "2025-08-14T22:01:00Z"
  },
  "timestamp": "2025-08-14T22:01:00Z"
}
```

**Error Response (404):**
```json
{
  "success": false,
  "error": "User not found",
  "timestamp": "2025-08-14T22:00:00Z"
}
```

### Delete User
**DELETE** `/api/v1/users/{id}`

Soft deletes a user (sets deleted_at timestamp).

**Parameters:**
- `id` (path, required): User ID

**Response (200):**
```json
{
  "success": true,
  "message": "User deleted successfully",
  "timestamp": "2025-08-14T22:00:00Z"
}
```

**Error Response (404):**
```json
{
  "success": false,
  "error": "User not found",
  "timestamp": "2025-08-14T22:00:00Z"
}
```

## cURL Examples

### Create a user with all fields:
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Smith",
    "email": "jane@example.com",
    "phone": "+1-555-0125",
    "address": "789 Pine Road"
  }'
```

### Create a user without address:
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bob Wilson",
    "email": "bob@example.com",
    "phone": "+1-555-0126"
  }'
```

### Get all users:
```bash
curl -X GET http://localhost:8080/api/v1/users
```

### Update a user:
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe Updated",
    "email": "john.doe@example.com",
    "phone": "+1-555-0127",
    "address": "999 New Street"
  }'
```

### Delete a user:
```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

## Status Codes

- `200 OK` - Successful GET, PUT, DELETE
- `201 Created` - Successful POST
- `400 Bad Request` - Invalid request body or validation error
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## Field Validation

### Name
- Required
- Must be non-empty string
- No specific format restrictions

### Email
- Required
- Must be valid email format
- Must be unique across all users

### Phone
- Required for API requests
- Stored as nullable in database for backward compatibility
- No format validation currently enforced
- Accepts any string format

### Address
- Optional
- Can be omitted from request
- Stored as nullable string in database
