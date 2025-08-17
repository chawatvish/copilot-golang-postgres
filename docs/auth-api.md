# Authentication API Documentation

This document describes the authentication endpoints available in the API.

## Table of Contents

- [Overview](#overview)
- [Authentication Flow](#authentication-flow)
- [Endpoints](#endpoints)
- [Request/Response Examples](#requestresponse-examples)
- [Error Handling](#error-handling)
- [JWT Token Structure](#jwt-token-structure)

## Overview

The API uses JWT (JSON Web Token) based authentication. Users can register, login, and perform authenticated operations using bearer tokens.

### Key Features

- User registration with email verification support
- User login with JWT token generation
- Password reset functionality
- Change password for authenticated users
- Token refresh
- User profile management
- Secure password hashing using bcrypt

## Authentication Flow

1. **Register**: Create a new user account
2. **Login**: Authenticate and receive JWT token
3. **Access Protected Routes**: Include JWT token in Authorization header
4. **Refresh Token**: Get a new token before expiry
5. **Logout**: Invalidate session (client-side token removal)

## Endpoints

### Public Endpoints (No Authentication Required)

#### POST /api/v1/auth/register

Register a new user account.

**Request Body:**

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "confirm_password": "password123",
  "phone": "+1-555-0101",
  "address": "123 Main St, New York, NY 10001"
}
```

**Response (201 Created):**

```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "phone": "+1-555-0101",
      "address": "123 Main St, New York, NY 10001",
      "is_active": true,
      "is_email_verified": false,
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z",
      "last_login_at": "2024-01-01T12:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 86400
  }
}
```

#### POST /api/v1/auth/login

Authenticate user and receive JWT token.

**Request Body:**

```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response (200 OK):**

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "phone": "+1-555-0101",
      "address": "123 Main St, New York, NY 10001",
      "is_active": true,
      "is_email_verified": true,
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z",
      "last_login_at": "2024-01-01T12:00:30Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 86400
  }
}
```

#### POST /api/v1/auth/forgot-password

Request password reset token.

**Request Body:**

```json
{
  "email": "john@example.com"
}
```

**Response (200 OK):**

```json
{
  "success": true,
  "message": "If the email exists, a password reset link has been sent",
  "data": null
}
```

#### POST /api/v1/auth/reset-password

Reset password using reset token.

**Request Body:**

```json
{
  "token": "uuid-reset-token-here",
  "new_password": "newpassword123",
  "confirm_password": "newpassword123"
}
```

**Response (200 OK):**

```json
{
  "success": true,
  "message": "Password reset successfully",
  "data": null
}
```

### Protected Endpoints (Authentication Required)

All protected endpoints require the `Authorization` header with a valid JWT token:

```
Authorization: Bearer <jwt_token>
```

#### GET /api/v1/auth/me

Get current authenticated user information.

**Response (200 OK):**

```json
{
  "success": true,
  "message": "User information retrieved",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+1-555-0101",
    "address": "123 Main St, New York, NY 10001",
    "is_active": true,
    "is_email_verified": true,
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z",
    "last_login_at": "2024-01-01T12:00:30Z"
  }
}
```

#### POST /api/v1/auth/change-password

Change password for authenticated user.

**Request Body:**

```json
{
  "current_password": "oldpassword123",
  "new_password": "newpassword123",
  "confirm_password": "newpassword123"
}
```

**Response (200 OK):**

```json
{
  "success": true,
  "message": "Password changed successfully",
  "data": null
}
```

#### POST /api/v1/auth/refresh-token

Refresh JWT token.

**Response (200 OK):**

```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": {
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com"
      // ... user data
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 86400
  }
}
```

#### POST /api/v1/auth/logout

Logout user (client should remove token).

**Response (200 OK):**

```json
{
  "success": true,
  "message": "Logged out successfully",
  "data": null
}
```

## Error Handling

All endpoints return consistent error responses:

**400 Bad Request:**

```json
{
  "success": false,
  "error": "Passwords do not match"
}
```

**401 Unauthorized:**

```json
{
  "success": false,
  "error": "Invalid email or password"
}
```

**422 Validation Error:**

```json
{
  "success": false,
  "error": "Key: 'RegisterRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"
}
```

**500 Internal Server Error:**

```json
{
  "success": false,
  "error": "Internal server error"
}
```

## JWT Token Structure

The JWT token contains the following claims:

```json
{
  "user_id": 1,
  "email": "john@example.com",
  "iat": 1704110400,
  "exp": 1704196800,
  "sub": "1"
}
```

### Token Usage

Include the token in the Authorization header for protected requests:

```bash
curl -H "Authorization: Bearer <your_jwt_token>" \
     -X GET http://localhost:8080/api/v1/auth/me
```

### Token Expiry

- Default expiry: 24 hours (configurable via `JWT_EXPIRE_HOUR`)
- Tokens should be refreshed before expiry using `/auth/refresh-token`
- Expired tokens will return 401 Unauthorized

## Configuration

JWT behavior can be configured using environment variables:

- `JWT_SECRET`: Secret key for signing tokens (should be at least 32 characters in production)
- `JWT_EXPIRE_HOUR`: Token expiry time in hours (default: 24)

## Security Considerations

1. **Password Security**: Passwords are hashed using bcrypt with default cost (10)
2. **Token Security**: Use HTTPS in production to prevent token interception
3. **Secret Key**: Use a strong, randomly generated JWT secret in production
4. **Token Storage**: Store tokens securely on the client side (consider httpOnly cookies for web apps)
5. **Password Reset**: Reset tokens expire in 1 hour for security
6. **Account Status**: Inactive accounts cannot login
7. **Email Verification**: Support for email verification workflow (tokens provided)
