# Setokin API Documentation

**Version:** 1.0.0  
**Base URL:** `https://api.setokin.com/v1`  
**Protocol:** HTTPS only

## Table of Contents

1. [Overview](#overview)
2. [Authentication](#authentication)
3. [Common Patterns](#common-patterns)
4. [Error Handling](#error-handling)
5. [API Endpoints](#api-endpoints)
   - [Authentication](#authentication-endpoints)
   - [Users](#users-endpoints)
   - [Categories](#categories-endpoints)
   - [Units](#units-endpoints)
   - [Items](#items-endpoints)
   - [Suppliers](#suppliers-endpoints)
   - [Batches](#batches-endpoints)
   - [Stock In](#stock-in-endpoints)
   - [Stock Out](#stock-out-endpoints)
   - [Inventory](#inventory-endpoints)
   - [Reports](#reports-endpoints)
   - [File Upload](#file-upload-endpoints)

---

## Overview

Setokin API adalah RESTful API untuk sistem manajemen inventory F&B dengan fitur FEFO (First Expired First Out).

### Key Features

- JWT-based authentication dengan dual token (access & refresh)
- FEFO automatic batch deduction
- Expiry tracking dan alerts
- Comprehensive reporting
- MinIO presigned URL untuk file uploads

### API Conventions

- All requests and responses use `application/json` content type
- Timestamps are in ISO 8601 format with timezone (RFC3339)
- UUIDs are used for all resource identifiers
- Pagination uses cursor-based approach
- All endpoints require HTTPS

---

## Authentication

Setokin menggunakan JWT (JSON Web Token) dengan dual token strategy:


### Token Types

| Token Type | Lifetime | Purpose | Storage |
|------------|----------|---------|---------|
| Access Token | 15 minutes | API authentication | Memory/localStorage |
| Refresh Token | 7 days | Renew access token | httpOnly cookie |

### Authentication Flow

```
1. User login → Receive access + refresh tokens
2. Use access token in Authorization header
3. Access token expires → Use refresh token to get new access token
4. Refresh token expires → User must login again
```

### Authorization Header

All authenticated endpoints require:

```
Authorization: Bearer <access_token>
```

### Token Refresh

When access token expires (401 with `token_expired` error), use refresh endpoint to obtain new tokens.

---

## Common Patterns

### Pagination

List endpoints support cursor-based pagination:

**Query Parameters:**
- `limit` (integer, default: 20, max: 100) - Number of items per page
- `cursor` (string, optional) - Cursor for next page

**Response:**
```json
{
  "data": [...],
  "pagination": {
    "next_cursor": "eyJpZCI6IjEyMyJ9",
    "has_more": true,
    "total": 150
  }
}
```

### Filtering

List endpoints support filtering via query parameters:

```
GET /items?category_id=<uuid>&is_active=true
```

### Sorting

Use `sort` parameter with field name and direction:

```
GET /items?sort=name:asc
GET /items?sort=created_at:desc
```

### Field Selection

Use `fields` parameter to select specific fields:

```
GET /items?fields=id,name,category_id
```

### Date Ranges

Use ISO 8601 format for date parameters:

```
GET /reports/daily?start_date=2024-03-01&end_date=2024-03-31
```

---

## Error Handling

### Error Response Format

```json
{
  "error": {
    "code": "validation_error",
    "message": "Validation failed",
    "details": [
      {
        "field": "email",
        "message": "Invalid email format"
      }
    ],
    "request_id": "req_abc123"
  }
}
```


### HTTP Status Codes

| Code | Meaning | Usage |
|------|---------|-------|
| 200 | OK | Successful GET, PUT, PATCH |
| 201 | Created | Successful POST |
| 204 | No Content | Successful DELETE |
| 400 | Bad Request | Invalid request format/parameters |
| 401 | Unauthorized | Missing or invalid authentication |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource not found |
| 409 | Conflict | Resource conflict (duplicate, constraint violation) |
| 422 | Unprocessable Entity | Validation error |
| 429 | Too Many Requests | Rate limit exceeded |
| 500 | Internal Server Error | Server error |
| 503 | Service Unavailable | Service temporarily unavailable |

### Error Codes

| Code | Description |
|------|-------------|
| `validation_error` | Request validation failed |
| `authentication_required` | Authentication required |
| `token_expired` | Access token expired |
| `token_invalid` | Invalid token |
| `insufficient_permissions` | User lacks required permissions |
| `resource_not_found` | Requested resource not found |
| `duplicate_resource` | Resource already exists |
| `insufficient_stock` | Not enough stock for operation |
| `batch_depleted` | Batch is already depleted |
| `invalid_quantity` | Invalid quantity value |
| `rate_limit_exceeded` | Too many requests |
| `internal_error` | Internal server error |

---

## API Endpoints

---

## Authentication Endpoints

### Register User

Create a new user account.

**Endpoint:** `POST /auth/register`

**Authentication:** None

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "full_name": "John Doe",
  "role": "staff"
}
```

**Request Fields:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| email | string | Yes | Valid email address |
| password | string | Yes | Min 8 chars, must contain uppercase, lowercase, number |
| full_name | string | Yes | User's full name |
| role | string | No | User role: `owner`, `manager`, `staff` (default: `staff`) |

**Success Response (201):**
```json
{
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "full_name": "John Doe",
      "role": "staff",
      "is_active": true,
      "created_at": "2024-03-11T10:30:00Z"
    },
    "tokens": {
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "expires_in": 900
    }
  }
}
```


**Error Responses:**

**400 Bad Request** - Invalid input
```json
{
  "error": {
    "code": "validation_error",
    "message": "Validation failed",
    "details": [
      {
        "field": "email",
        "message": "Invalid email format"
      },
      {
        "field": "password",
        "message": "Password must be at least 8 characters"
      }
    ]
  }
}
```

**409 Conflict** - Email already exists
```json
{
  "error": {
    "code": "duplicate_resource",
    "message": "Email already registered"
  }
}
```

---

### Login

Authenticate user and receive tokens.

**Endpoint:** `POST /auth/login`

**Authentication:** None

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**Success Response (200):**
```json
{
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "full_name": "John Doe",
      "role": "staff",
      "is_active": true
    },
    "tokens": {
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "expires_in": 900
    }
  }
}
```

**Error Responses:**

**401 Unauthorized** - Invalid credentials
```json
{
  "error": {
    "code": "authentication_required",
    "message": "Invalid email or password"
  }
}
```

**403 Forbidden** - Account inactive
```json
{
  "error": {
    "code": "insufficient_permissions",
    "message": "Account is inactive"
  }
}
```

---


### Refresh Token

Get new access token using refresh token.

**Endpoint:** `POST /auth/refresh`

**Authentication:** Refresh token required

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Success Response (200):**
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 900
  }
}
```

**Error Responses:**

**401 Unauthorized** - Invalid or expired refresh token
```json
{
  "error": {
    "code": "token_invalid",
    "message": "Invalid or expired refresh token"
  }
}
```

---

### Logout

Revoke refresh token and invalidate session.

**Endpoint:** `POST /auth/logout`

**Authentication:** Required

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Success Response (204):**
No content

---

### Get Current User

Get authenticated user information.

**Endpoint:** `GET /auth/me`

**Authentication:** Required

**Success Response (200):**
```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "full_name": "John Doe",
    "role": "staff",
    "is_active": true,
    "created_at": "2024-03-11T10:30:00Z",
    "updated_at": "2024-03-11T10:30:00Z"
  }
}
```

---


## Categories Endpoints

### List Categories

Get all categories.

**Endpoint:** `GET /categories`

**Authentication:** Required

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| limit | integer | No | Items per page (default: 20, max: 100) |
| cursor | string | No | Pagination cursor |
| sort | string | No | Sort field and direction (e.g., `name:asc`) |

**Success Response (200):**
```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Daging",
      "description": "Daging sapi, ayam, ikan, dll",
      "created_at": "2024-03-11T10:30:00Z",
      "updated_at": "2024-03-11T10:30:00Z"
    },
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "name": "Sayuran",
      "description": "Sayuran segar",
      "created_at": "2024-03-11T10:30:00Z",
      "updated_at": "2024-03-11T10:30:00Z"
    }
  ],
  "pagination": {
    "next_cursor": null,
    "has_more": false,
    "total": 2
  }
}
```

---

### Get Category

Get single category by ID.

**Endpoint:** `GET /categories/{id}`

**Authentication:** Required

**Path Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | uuid | Yes | Category ID |

**Success Response (200):**
```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Daging",
    "description": "Daging sapi, ayam, ikan, dll",
    "created_at": "2024-03-11T10:30:00Z",
    "updated_at": "2024-03-11T10:30:00Z"
  }
}
```

**Error Responses:**

**404 Not Found**
```json
{
  "error": {
    "code": "resource_not_found",
    "message": "Category not found"
  }
}
```

---


### Create Category

Create a new category.

**Endpoint:** `POST /categories`

**Authentication:** Required (role: owner, manager)

**Request Body:**
```json
{
  "name": "Bumbu",
  "description": "Bumbu dapur dan rempah"
}
```

**Success Response (201):**
```json
{
  "data": {
    "id": "770e8400-e29b-41d4-a716-446655440002",
    "name": "Bumbu",
    "description": "Bumbu dapur dan rempah",
    "created_at": "2024-03-11T10:30:00Z",
    "updated_at": "2024-03-11T10:30:00Z"
  }
}
```

**Error Responses:**

**409 Conflict** - Duplicate name
```json
{
  "error": {
    "code": "duplicate_resource",
    "message": "Category name already exists"
  }
}
```

---

### Update Category

Update existing category.

**Endpoint:** `PUT /categories/{id}`

**Authentication:** Required (role: owner, manager)

**Request Body:**
```json
{
  "name": "Bumbu & Rempah",
  "description": "Bumbu dapur, rempah, dan seasoning"
}
```

**Success Response (200):**
```json
{
  "data": {
    "id": "770e8400-e29b-41d4-a716-446655440002",
    "name": "Bumbu & Rempah",
    "description": "Bumbu dapur, rempah, dan seasoning",
    "created_at": "2024-03-11T10:30:00Z",
    "updated_at": "2024-03-11T11:45:00Z"
  }
}
```

---

### Delete Category

Delete a category.

**Endpoint:** `DELETE /categories/{id}`

**Authentication:** Required (role: owner, manager)

**Success Response (204):**
No content

**Error Responses:**

**409 Conflict** - Category in use
```json
{
  "error": {
    "code": "resource_conflict",
    "message": "Cannot delete category with existing items"
  }
}
```

---


## Units Endpoints

### List Units

Get all units of measurement.

**Endpoint:** `GET /units`

**Authentication:** Required

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| type | string | No | Filter by type: `weight`, `volume`, `count` |

**Success Response (200):**
```json
{
  "data": [
    {
      "id": "880e8400-e29b-41d4-a716-446655440003",
      "name": "Kilogram",
      "abbreviation": "kg",
      "type": "weight",
      "created_at": "2024-03-11T10:30:00Z"
    },
    {
      "id": "990e8400-e29b-41d4-a716-446655440004",
      "name": "Liter",
      "abbreviation": "L",
      "type": "volume",
      "created_at": "2024-03-11T10:30:00Z"
    }
  ],
  "pagination": {
    "next_cursor": null,
    "has_more": false,
    "total": 2
  }
}
```

---

## Items Endpoints

### List Items

Get all inventory items.

**Endpoint:** `GET /items`

**Authentication:** Required

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| limit | integer | No | Items per page (default: 20, max: 100) |
| cursor | string | No | Pagination cursor |
| category_id | uuid | No | Filter by category |
| is_active | boolean | No | Filter by active status |
| search | string | No | Search by name |
| sort | string | No | Sort (e.g., `name:asc`, `created_at:desc`) |

**Success Response (200):**
```json
{
  "data": [
    {
      "id": "aa0e8400-e29b-41d4-a716-446655440005",
      "name": "Ayam Fillet",
      "category": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "Daging"
      },
      "unit": {
        "id": "880e8400-e29b-41d4-a716-446655440003",
        "abbreviation": "kg"
      },
      "minimum_stock": 5.000,
      "description": "Ayam fillet segar",
      "is_active": true,
      "created_at": "2024-03-11T10:30:00Z",
      "updated_at": "2024-03-11T10:30:00Z"
    }
  ],
  "pagination": {
    "next_cursor": "eyJpZCI6ImFhMGU4NDAwIn0",
    "has_more": true,
    "total": 45
  }
}
```

---


### Get Item

Get single item by ID.

**Endpoint:** `GET /items/{id}`

**Authentication:** Required

**Success Response (200):**
```json
{
  "data": {
    "id": "aa0e8400-e29b-41d4-a716-446655440005",
    "name": "Ayam Fillet",
    "category": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Daging"
    },
    "unit": {
      "id": "880e8400-e29b-41d4-a716-446655440003",
      "name": "Kilogram",
      "abbreviation": "kg"
    },
    "minimum_stock": 5.000,
    "description": "Ayam fillet segar",
    "is_active": true,
    "current_stock": 12.500,
    "active_batches": 2,
    "is_low_stock": false,
    "created_at": "2024-03-11T10:30:00Z",
    "updated_at": "2024-03-11T10:30:00Z"
  }
}
```

---

### Create Item

Create a new inventory item.

**Endpoint:** `POST /items`

**Authentication:** Required (role: owner, manager)

**Request Body:**
```json
{
  "name": "Tepung Terigu",
  "category_id": "770e8400-e29b-41d4-a716-446655440002",
  "unit_id": "880e8400-e29b-41d4-a716-446655440003",
  "minimum_stock": 10.000,
  "description": "Tepung terigu protein sedang"
}
```

**Request Fields:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | Yes | Item name |
| category_id | uuid | Yes | Category ID |
| unit_id | uuid | Yes | Unit ID |
| minimum_stock | number | No | Minimum stock threshold (default: 0) |
| description | string | No | Item description |

**Success Response (201):**
```json
{
  "data": {
    "id": "bb0e8400-e29b-41d4-a716-446655440006",
    "name": "Tepung Terigu",
    "category": {
      "id": "770e8400-e29b-41d4-a716-446655440002",
      "name": "Bumbu & Rempah"
    },
    "unit": {
      "id": "880e8400-e29b-41d4-a716-446655440003",
      "abbreviation": "kg"
    },
    "minimum_stock": 10.000,
    "description": "Tepung terigu protein sedang",
    "is_active": true,
    "created_at": "2024-03-11T10:30:00Z",
    "updated_at": "2024-03-11T10:30:00Z"
  }
}
```

---


### Update Item

Update existing item.

**Endpoint:** `PUT /items/{id}`

**Authentication:** Required (role: owner, manager)

**Request Body:**
```json
{
  "name": "Tepung Terigu Premium",
  "category_id": "770e8400-e29b-41d4-a716-446655440002",
  "unit_id": "880e8400-e29b-41d4-a716-446655440003",
  "minimum_stock": 15.000,
  "description": "Tepung terigu protein tinggi",
  "is_active": true
}
```

**Success Response (200):**
```json
{
  "data": {
    "id": "bb0e8400-e29b-41d4-a716-446655440006",
    "name": "Tepung Terigu Premium",
    "category": {
      "id": "770e8400-e29b-41d4-a716-446655440002",
      "name": "Bumbu & Rempah"
    },
    "unit": {
      "id": "880e8400-e29b-41d4-a716-446655440003",
      "abbreviation": "kg"
    },
    "minimum_stock": 15.000,
    "description": "Tepung terigu protein tinggi",
    "is_active": true,
    "created_at": "2024-03-11T10:30:00Z",
    "updated_at": "2024-03-11T11:45:00Z"
  }
}
```

---

### Delete Item

Delete an item (soft delete - sets is_active to false).

**Endpoint:** `DELETE /items/{id}`

**Authentication:** Required (role: owner, manager)

**Success Response (204):**
No content

**Error Responses:**

**409 Conflict** - Item has active batches
```json
{
  "error": {
    "code": "resource_conflict",
    "message": "Cannot delete item with active batches"
  }
}
```

---


## Suppliers Endpoints

### List Suppliers

Get all suppliers.

**Endpoint:** `GET /suppliers`

**Authentication:** Required

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| limit | integer | No | Items per page (default: 20, max: 100) |
| cursor | string | No | Pagination cursor |
| is_active | boolean | No | Filter by active status |
| search | string | No | Search by name |

**Success Response (200):**
```json
{
  "data": [
    {
      "id": "cc0e8400-e29b-41d4-a716-446655440007",
      "name": "PT Sumber Pangan",
      "contact_person": "Budi Santoso",
      "phone": "+62812345678",
      "email": "budi@sumberpangan.com",
      "address": "Jl. Raya No. 123, Jakarta",
      "is_active": true,
      "created_at": "2024-03-11T10:30:00Z",
      "updated_at": "2024-03-11T10:30:00Z"
    }
  ],
  "pagination": {
    "next_cursor": null,
    "has_more": false,
    "total": 1
  }
}
```

---

### Create Supplier

Create a new supplier.

**Endpoint:** `POST /suppliers`

**Authentication:** Required (role: owner, manager)

**Request Body:**
```json
{
  "name": "CV Mitra Dagang",
  "contact_person": "Siti Aminah",
  "phone": "+62823456789",
  "email": "siti@mitradagang.com",
  "address": "Jl. Perdagangan No. 45, Bandung"
}
```

**Success Response (201):**
```json
{
  "data": {
    "id": "dd0e8400-e29b-41d4-a716-446655440008",
    "name": "CV Mitra Dagang",
    "contact_person": "Siti Aminah",
    "phone": "+62823456789",
    "email": "siti@mitradagang.com",
    "address": "Jl. Perdagangan No. 45, Bandung",
    "is_active": true,
    "created_at": "2024-03-11T10:30:00Z",
    "updated_at": "2024-03-11T10:30:00Z"
  }
}
```

---


## Stock In Endpoints

### Create Stock In

Record incoming stock (purchase).

**Endpoint:** `POST /stock-in`

**Authentication:** Required

**Request Body:**
```json
{
  "item_id": "aa0e8400-e29b-41d4-a716-446655440005",
  "quantity": 10.000,
  "purchase_date": "2024-03-11",
  "expiry_date": "2024-03-18",
  "supplier_id": "cc0e8400-e29b-41d4-a716-446655440007",
  "purchase_price": 150000.00,
  "notes": "Pembelian rutin mingguan"
}
```

**Request Fields:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| item_id | uuid | Yes | Item ID |
| quantity | number | Yes | Quantity (must be > 0) |
| purchase_date | date | Yes | Purchase date (YYYY-MM-DD) |
| expiry_date | date | Yes | Expiry date (YYYY-MM-DD) |
| supplier_id | uuid | No | Supplier ID |
| purchase_price | number | No | Purchase price |
| notes | string | No | Additional notes |

**Success Response (201):**
```json
{
  "data": {
    "stock_in": {
      "id": "ee0e8400-e29b-41d4-a716-446655440009",
      "item": {
        "id": "aa0e8400-e29b-41d4-a716-446655440005",
        "name": "Ayam Fillet"
      },
      "batch": {
        "id": "ff0e8400-e29b-41d4-a716-44665544000a",
        "batch_number": "BATCH-20240311-0001"
      },
      "quantity": 10.000,
      "purchase_date": "2024-03-11",
      "supplier": {
        "id": "cc0e8400-e29b-41d4-a716-446655440007",
        "name": "PT Sumber Pangan"
      },
      "purchase_price": 150000.00,
      "notes": "Pembelian rutin mingguan",
      "created_by": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "full_name": "John Doe"
      },
      "created_at": "2024-03-11T10:30:00Z"
    },
    "batch": {
      "id": "ff0e8400-e29b-41d4-a716-44665544000a",
      "batch_number": "BATCH-20240311-0001",
      "initial_quantity": 10.000,
      "remaining_quantity": 10.000,
      "expiry_date": "2024-03-18",
      "is_depleted": false
    }
  }
}
```

**Error Responses:**

**422 Unprocessable Entity** - Invalid quantity
```json
{
  "error": {
    "code": "invalid_quantity",
    "message": "Quantity must be greater than 0"
  }
}
```

**422 Unprocessable Entity** - Invalid expiry date
```json
{
  "error": {
    "code": "validation_error",
    "message": "Expiry date must be in the future"
  }
}
```

---


### List Stock In

Get stock in history.

**Endpoint:** `GET /stock-in`

**Authentication:** Required

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| limit | integer | No | Items per page (default: 20, max: 100) |
| cursor | string | No | Pagination cursor |
| item_id | uuid | No | Filter by item |
| supplier_id | uuid | No | Filter by supplier |
| start_date | date | No | Filter from date (YYYY-MM-DD) |
| end_date | date | No | Filter to date (YYYY-MM-DD) |
| sort | string | No | Sort (default: `created_at:desc`) |

**Success Response (200):**
```json
{
  "data": [
    {
      "id": "ee0e8400-e29b-41d4-a716-446655440009",
      "item": {
        "id": "aa0e8400-e29b-41d4-a716-446655440005",
        "name": "Ayam Fillet"
      },
      "batch": {
        "id": "ff0e8400-e29b-41d4-a716-44665544000a",
        "batch_number": "BATCH-20240311-0001"
      },
      "quantity": 10.000,
      "purchase_date": "2024-03-11",
      "supplier": {
        "id": "cc0e8400-e29b-41d4-a716-446655440007",
        "name": "PT Sumber Pangan"
      },
      "purchase_price": 150000.00,
      "notes": "Pembelian rutin mingguan",
      "created_by": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "full_name": "John Doe"
      },
      "created_at": "2024-03-11T10:30:00Z"
    }
  ],
  "pagination": {
    "next_cursor": "eyJpZCI6ImVlMGU4NDAwIn0",
    "has_more": true,
    "total": 125
  }
}
```

---


## Stock Out Endpoints

### Create Stock Out

Record stock usage with automatic FEFO deduction.

**Endpoint:** `POST /stock-out`

**Authentication:** Required

**Request Body:**
```json
{
  "item_id": "aa0e8400-e29b-41d4-a716-446655440005",
  "quantity": 3.000,
  "usage_date": "2024-03-12",
  "notes": "Dipakai untuk menu ayam goreng"
}
```

**Request Fields:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| item_id | uuid | Yes | Item ID |
| quantity | number | Yes | Quantity to use (must be > 0) |
| usage_date | date | Yes | Usage date (YYYY-MM-DD) |
| notes | string | No | Usage notes |

**Success Response (201):**
```json
{
  "data": {
    "stock_out": {
      "id": "110e8400-e29b-41d4-a716-44665544000b",
      "item": {
        "id": "aa0e8400-e29b-41d4-a716-446655440005",
        "name": "Ayam Fillet"
      },
      "quantity": 3.000,
      "usage_date": "2024-03-12",
      "notes": "Dipakai untuk menu ayam goreng",
      "created_by": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "full_name": "John Doe"
      },
      "created_at": "2024-03-12T08:15:00Z"
    },
    "deductions": [
      {
        "batch": {
          "id": "ff0e8400-e29b-41d4-a716-44665544000a",
          "batch_number": "BATCH-20240311-0001",
          "expiry_date": "2024-03-15"
        },
        "quantity_used": 3.000,
        "remaining_after": 7.000
      }
    ],
    "remaining_stock": 7.000
  }
}
```

**FEFO Logic Example:**

If there are multiple batches:
```json
{
  "deductions": [
    {
      "batch": {
        "id": "batch-1",
        "batch_number": "BATCH-20240310-0001",
        "expiry_date": "2024-03-15"
      },
      "quantity_used": 2.000,
      "remaining_after": 0.000
    },
    {
      "batch": {
        "id": "batch-2",
        "batch_number": "BATCH-20240311-0001",
        "expiry_date": "2024-03-20"
      },
      "quantity_used": 1.000,
      "remaining_after": 9.000
    }
  ]
}
```

**Error Responses:**

**422 Unprocessable Entity** - Insufficient stock
```json
{
  "error": {
    "code": "insufficient_stock",
    "message": "Insufficient stock. Available: 2.5 kg, Requested: 3.0 kg"
  }
}
```

---


### List Stock Out

Get stock out history.

**Endpoint:** `GET /stock-out`

**Authentication:** Required

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| limit | integer | No | Items per page (default: 20, max: 100) |
| cursor | string | No | Pagination cursor |
| item_id | uuid | No | Filter by item |
| start_date | date | No | Filter from date (YYYY-MM-DD) |
| end_date | date | No | Filter to date (YYYY-MM-DD) |
| sort | string | No | Sort (default: `created_at:desc`) |

**Success Response (200):**
```json
{
  "data": [
    {
      "id": "110e8400-e29b-41d4-a716-44665544000b",
      "item": {
        "id": "aa0e8400-e29b-41d4-a716-446655440005",
        "name": "Ayam Fillet"
      },
      "quantity": 3.000,
      "usage_date": "2024-03-12",
      "notes": "Dipakai untuk menu ayam goreng",
      "created_by": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "full_name": "John Doe"
      },
      "created_at": "2024-03-12T08:15:00Z"
    }
  ],
  "pagination": {
    "next_cursor": "eyJpZCI6IjExMGU4NDAwIn0",
    "has_more": true,
    "total": 89
  }
}
```

---

### Get Stock Out Details

Get detailed batch deductions for a stock out transaction.

**Endpoint:** `GET /stock-out/{id}/details`

**Authentication:** Required

**Success Response (200):**
```json
{
  "data": {
    "stock_out": {
      "id": "110e8400-e29b-41d4-a716-44665544000b",
      "item": {
        "id": "aa0e8400-e29b-41d4-a716-446655440005",
        "name": "Ayam Fillet"
      },
      "quantity": 3.000,
      "usage_date": "2024-03-12"
    },
    "deductions": [
      {
        "id": "120e8400-e29b-41d4-a716-44665544000c",
        "batch": {
          "id": "ff0e8400-e29b-41d4-a716-44665544000a",
          "batch_number": "BATCH-20240311-0001",
          "expiry_date": "2024-03-15"
        },
        "quantity_used": 3.000,
        "created_at": "2024-03-12T08:15:00Z"
      }
    ]
  }
}
```

---


## Batches Endpoints

### List Batches

Get all batches for an item.

**Endpoint:** `GET /items/{item_id}/batches`

**Authentication:** Required

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| is_depleted | boolean | No | Filter by depletion status |
| sort | string | No | Sort (default: `expiry_date:asc`) |

**Success Response (200):**
```json
{
  "data": [
    {
      "id": "ff0e8400-e29b-41d4-a716-44665544000a",
      "batch_number": "BATCH-20240311-0001",
      "initial_quantity": 10.000,
      "remaining_quantity": 7.000,
      "expiry_date": "2024-03-15",
      "is_depleted": false,
      "status": "expiring_soon",
      "days_until_expiry": 3,
      "created_at": "2024-03-11T10:30:00Z",
      "updated_at": "2024-03-12T08:15:00Z"
    },
    {
      "id": "130e8400-e29b-41d4-a716-44665544000d",
      "batch_number": "BATCH-20240312-0002",
      "initial_quantity": 5.000,
      "remaining_quantity": 5.000,
      "expiry_date": "2024-03-20",
      "is_depleted": false,
      "status": "good",
      "days_until_expiry": 8,
      "created_at": "2024-03-12T09:00:00Z",
      "updated_at": "2024-03-12T09:00:00Z"
    }
  ]
}
```

**Batch Status Values:**
- `expired` - Expiry date has passed
- `expiring_soon` - Expires within 3 days
- `good` - More than 3 days until expiry

---

### Get Batch

Get single batch details.

**Endpoint:** `GET /batches/{id}`

**Authentication:** Required

**Success Response (200):**
```json
{
  "data": {
    "id": "ff0e8400-e29b-41d4-a716-44665544000a",
    "batch_number": "BATCH-20240311-0001",
    "item": {
      "id": "aa0e8400-e29b-41d4-a716-446655440005",
      "name": "Ayam Fillet",
      "unit": "kg"
    },
    "initial_quantity": 10.000,
    "remaining_quantity": 7.000,
    "expiry_date": "2024-03-15",
    "is_depleted": false,
    "status": "expiring_soon",
    "days_until_expiry": 3,
    "usage_history": [
      {
        "stock_out_id": "110e8400-e29b-41d4-a716-44665544000b",
        "quantity_used": 3.000,
        "usage_date": "2024-03-12",
        "created_at": "2024-03-12T08:15:00Z"
      }
    ],
    "created_at": "2024-03-11T10:30:00Z",
    "updated_at": "2024-03-12T08:15:00Z"
  }
}
```

---


## Inventory Endpoints

### Get Current Inventory

Get current stock levels for all items.

**Endpoint:** `GET /inventory`

**Authentication:** Required

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| limit | integer | No | Items per page (default: 20, max: 100) |
| cursor | string | No | Pagination cursor |
| category_id | uuid | No | Filter by category |
| low_stock_only | boolean | No | Show only low stock items |
| search | string | No | Search by item name |

**Success Response (200):**
```json
{
  "data": [
    {
      "item_id": "aa0e8400-e29b-41d4-a716-446655440005",
      "item_name": "Ayam Fillet",
      "category_name": "Daging",
      "unit": "kg",
      "total_stock": 12.000,
      "minimum_stock": 5.000,
      "is_low_stock": false,
      "active_batches": 2,
      "stock_status": "adequate"
    },
    {
      "item_id": "bb0e8400-e29b-41d4-a716-446655440006",
      "item_name": "Tepung Terigu",
      "category_name": "Bumbu & Rempah",
      "unit": "kg",
      "total_stock": 3.500,
      "minimum_stock": 10.000,
      "is_low_stock": true,
      "active_batches": 1,
      "stock_status": "low"
    }
  ],
  "pagination": {
    "next_cursor": "eyJpZCI6ImJiMGU4NDAwIn0",
    "has_more": true,
    "total": 45
  },
  "summary": {
    "total_items": 45,
    "low_stock_items": 8,
    "out_of_stock_items": 2
  }
}
```

**Stock Status Values:**
- `out_of_stock` - Total stock = 0
- `low` - Total stock <= minimum stock
- `adequate` - Total stock > minimum stock

---

### Get Expiring Items

Get items with batches expiring soon (within 3 days).

**Endpoint:** `GET /inventory/expiring`

**Authentication:** Required

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| days | integer | No | Days threshold (default: 3, max: 30) |

**Success Response (200):**
```json
{
  "data": [
    {
      "batch_id": "ff0e8400-e29b-41d4-a716-44665544000a",
      "item_id": "aa0e8400-e29b-41d4-a716-446655440005",
      "item_name": "Ayam Fillet",
      "category_name": "Daging",
      "batch_number": "BATCH-20240311-0001",
      "remaining_quantity": 7.000,
      "unit": "kg",
      "expiry_date": "2024-03-15",
      "days_until_expiry": 3,
      "urgency": "high"
    },
    {
      "batch_id": "140e8400-e29b-41d4-a716-44665544000e",
      "item_id": "150e8400-e29b-41d4-a716-44665544000f",
      "item_name": "Susu Segar",
      "category_name": "Dairy",
      "batch_number": "BATCH-20240312-0003",
      "remaining_quantity": 5.000,
      "unit": "L",
      "expiry_date": "2024-03-13",
      "days_until_expiry": 1,
      "urgency": "critical"
    }
  ],
  "summary": {
    "total_expiring_batches": 2,
    "critical_count": 1,
    "high_count": 1,
    "total_value_at_risk": 275000.00
  }
}
```

**Urgency Levels:**
- `critical` - Expires today or tomorrow (0-1 days)
- `high` - Expires within 2-3 days
- `medium` - Expires within 4-7 days

---


## Reports Endpoints

### Daily Report

Get daily stock activity report.

**Endpoint:** `GET /reports/daily`

**Authentication:** Required

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| date | date | No | Report date (default: today, YYYY-MM-DD) |

**Success Response (200):**
```json
{
  "data": {
    "date": "2024-03-12",
    "stock_in": [
      {
        "item_name": "Ayam Fillet",
        "quantity": 10.000,
        "unit": "kg",
        "supplier": "PT Sumber Pangan",
        "purchase_price": 150000.00,
        "time": "10:30:00"
      }
    ],
    "stock_out": [
      {
        "item_name": "Ayam Fillet",
        "quantity": 3.000,
        "unit": "kg",
        "notes": "Dipakai untuk menu ayam goreng",
        "time": "08:15:00"
      },
      {
        "item_name": "Tepung Terigu",
        "quantity": 1.000,
        "unit": "kg",
        "notes": "Menu roti",
        "time": "09:30:00"
      }
    ],
    "summary": {
      "total_stock_in": 10.000,
      "total_stock_out": 4.000,
      "total_purchase_value": 150000.00,
      "transactions_count": 3
    }
  }
}
```

---

### Weekly Usage Report

Get weekly stock usage summary.

**Endpoint:** `GET /reports/weekly`

**Authentication:** Required

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| start_date | date | No | Week start date (default: current week Monday) |

**Success Response (200):**
```json
{
  "data": {
    "period": {
      "start_date": "2024-03-11",
      "end_date": "2024-03-17",
      "week_number": 11
    },
    "usage": [
      {
        "item_id": "aa0e8400-e29b-41d4-a716-446655440005",
        "item_name": "Ayam Fillet",
        "category": "Daging",
        "total_used": 18.000,
        "unit": "kg",
        "transaction_count": 6,
        "average_daily": 3.000
      },
      {
        "item_id": "bb0e8400-e29b-41d4-a716-446655440006",
        "item_name": "Tepung Terigu",
        "category": "Bumbu & Rempah",
        "total_used": 9.000,
        "unit": "kg",
        "transaction_count": 4,
        "average_daily": 1.500
      }
    ],
    "summary": {
      "total_items": 12,
      "total_transactions": 45,
      "most_used_item": "Ayam Fillet",
      "least_used_item": "Garam"
    }
  }
}
```

---


### Monthly Usage Report

Get monthly stock usage summary.

**Endpoint:** `GET /reports/monthly`

**Authentication:** Required

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| year | integer | No | Year (default: current year) |
| month | integer | No | Month 1-12 (default: current month) |

**Success Response (200):**
```json
{
  "data": {
    "period": {
      "year": 2024,
      "month": 3,
      "month_name": "March",
      "start_date": "2024-03-01",
      "end_date": "2024-03-31"
    },
    "usage": [
      {
        "item_id": "aa0e8400-e29b-41d4-a716-446655440005",
        "item_name": "Ayam Fillet",
        "category": "Daging",
        "total_used": 72.000,
        "unit": "kg",
        "transaction_count": 24,
        "average_daily": 2.400,
        "total_cost": 1080000.00
      },
      {
        "item_id": "bb0e8400-e29b-41d4-a716-446655440006",
        "item_name": "Tepung Terigu",
        "category": "Bumbu & Rempah",
        "total_used": 40.000,
        "unit": "kg",
        "transaction_count": 18,
        "average_daily": 1.333,
        "total_cost": 400000.00
      }
    ],
    "summary": {
      "total_items": 25,
      "total_transactions": 156,
      "total_cost": 5250000.00,
      "most_used_category": "Daging",
      "highest_cost_item": "Ayam Fillet"
    },
    "trends": {
      "vs_previous_month": {
        "usage_change_percent": 12.5,
        "cost_change_percent": 8.3,
        "trend": "increasing"
      }
    }
  }
}
```

---

### Stock Movement Report

Get detailed stock movement for an item.

**Endpoint:** `GET /reports/stock-movement/{item_id}`

**Authentication:** Required

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| start_date | date | Yes | Start date (YYYY-MM-DD) |
| end_date | date | Yes | End date (YYYY-MM-DD) |

**Success Response (200):**
```json
{
  "data": {
    "item": {
      "id": "aa0e8400-e29b-41d4-a716-446655440005",
      "name": "Ayam Fillet",
      "unit": "kg"
    },
    "period": {
      "start_date": "2024-03-01",
      "end_date": "2024-03-31"
    },
    "opening_stock": 5.000,
    "movements": [
      {
        "date": "2024-03-11",
        "type": "stock_in",
        "quantity": 10.000,
        "balance": 15.000,
        "reference": "BATCH-20240311-0001",
        "notes": "Pembelian rutin"
      },
      {
        "date": "2024-03-12",
        "type": "stock_out",
        "quantity": -3.000,
        "balance": 12.000,
        "reference": "SO-20240312-0001",
        "notes": "Menu ayam goreng"
      }
    ],
    "closing_stock": 12.000,
    "summary": {
      "total_in": 45.000,
      "total_out": 38.000,
      "net_change": 7.000
    }
  }
}
```

---


## File Upload Endpoints

Setokin uses MinIO for file storage with presigned URL strategy for secure uploads.

### Request Upload URL

Get presigned URL for file upload.

**Endpoint:** `POST /uploads/request`

**Authentication:** Required

**Request Body:**
```json
{
  "file_name": "invoice-2024-03-11.pdf",
  "file_type": "application/pdf",
  "file_size": 524288,
  "purpose": "stock_in_invoice"
}
```

**Request Fields:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| file_name | string | Yes | Original file name |
| file_type | string | Yes | MIME type |
| file_size | integer | Yes | File size in bytes (max: 10MB) |
| purpose | string | Yes | Upload purpose: `stock_in_invoice`, `item_image`, `supplier_document` |

**Success Response (200):**
```json
{
  "data": {
    "upload_id": "160e8400-e29b-41d4-a716-446655440010",
    "presigned_url": "https://minio.setokin.com/uploads/550e8400-e29b-41d4-a716-446655440000/invoice-2024-03-11.pdf?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=...",
    "file_key": "uploads/550e8400-e29b-41d4-a716-446655440000/invoice-2024-03-11.pdf",
    "expires_at": "2024-03-11T11:00:00Z",
    "max_file_size": 10485760
  }
}
```

**Upload Instructions:**

1. Use the `presigned_url` to upload file directly to MinIO
2. Make PUT request to presigned URL with file as body
3. Set `Content-Type` header to match `file_type`
4. Upload must complete before `expires_at`

**Example Upload (cURL):**
```bash
curl -X PUT \
  -H "Content-Type: application/pdf" \
  --data-binary @invoice-2024-03-11.pdf \
  "https://minio.setokin.com/uploads/..."
```

**Error Responses:**

**422 Unprocessable Entity** - File too large
```json
{
  "error": {
    "code": "validation_error",
    "message": "File size exceeds maximum allowed (10MB)"
  }
}
```

**422 Unprocessable Entity** - Invalid file type
```json
{
  "error": {
    "code": "validation_error",
    "message": "File type not allowed. Allowed types: image/jpeg, image/png, application/pdf"
  }
}
```

---

### Confirm Upload

Confirm successful file upload.

**Endpoint:** `POST /uploads/{upload_id}/confirm`

**Authentication:** Required

**Request Body:**
```json
{
  "file_key": "uploads/550e8400-e29b-41d4-a716-446655440000/invoice-2024-03-11.pdf"
}
```

**Success Response (200):**
```json
{
  "data": {
    "upload_id": "160e8400-e29b-41d4-a716-446655440010",
    "file_url": "https://cdn.setokin.com/uploads/550e8400-e29b-41d4-a716-446655440000/invoice-2024-03-11.pdf",
    "file_key": "uploads/550e8400-e29b-41d4-a716-446655440000/invoice-2024-03-11.pdf",
    "status": "confirmed",
    "confirmed_at": "2024-03-11T10:45:00Z"
  }
}
```

---


### Get Download URL

Get presigned URL for file download.

**Endpoint:** `GET /uploads/{upload_id}/download`

**Authentication:** Required

**Success Response (200):**
```json
{
  "data": {
    "download_url": "https://minio.setokin.com/uploads/550e8400-e29b-41d4-a716-446655440000/invoice-2024-03-11.pdf?X-Amz-Algorithm=AWS4-HMAC-SHA256&...",
    "expires_at": "2024-03-11T11:00:00Z",
    "file_name": "invoice-2024-03-11.pdf",
    "file_size": 524288,
    "file_type": "application/pdf"
  }
}
```

---

## Webhooks (Future)

Setokin will support webhooks for real-time notifications.

### Webhook Events

| Event | Description |
|-------|-------------|
| `stock.low` | Item stock falls below minimum |
| `batch.expiring` | Batch will expire within threshold |
| `batch.expired` | Batch has expired |
| `stock_in.created` | New stock in transaction |
| `stock_out.created` | New stock out transaction |

### Webhook Payload Example

```json
{
  "event": "stock.low",
  "timestamp": "2024-03-11T10:30:00Z",
  "data": {
    "item_id": "aa0e8400-e29b-41d4-a716-446655440005",
    "item_name": "Ayam Fillet",
    "current_stock": 3.500,
    "minimum_stock": 5.000,
    "unit": "kg"
  }
}
```

---

## Rate Limiting

API requests are rate limited to ensure fair usage.

### Limits

| Tier | Requests per Minute | Requests per Hour |
|------|---------------------|-------------------|
| Free | 60 | 1000 |
| Pro | 300 | 10000 |
| Enterprise | Custom | Custom |

### Rate Limit Headers

```
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 45
X-RateLimit-Reset: 1710158400
```

### Rate Limit Exceeded Response (429)

```json
{
  "error": {
    "code": "rate_limit_exceeded",
    "message": "Rate limit exceeded. Try again in 30 seconds.",
    "retry_after": 30
  }
}
```

---


## API Versioning

Setokin API uses URL-based versioning.

### Current Version

`v1` - Current stable version

### Version Format

```
https://api.setokin.com/v1/...
```

### Deprecation Policy

- Minimum 6 months notice before deprecation
- Deprecated endpoints return `Deprecation` header
- Migration guide provided for breaking changes

### Deprecation Header Example

```
Deprecation: true
Sunset: Sat, 31 Dec 2024 23:59:59 GMT
Link: <https://docs.setokin.com/migration/v2>; rel="deprecation"
```

---

## SDK & Client Libraries

### Official SDKs

- **JavaScript/TypeScript** - `npm install @setokin/sdk`
- **Go** - `go get github.com/setokin/go-sdk`
- **Python** - `pip install setokin-sdk` (Coming soon)

### SDK Example (TypeScript)

```typescript
import { SetokinClient } from '@setokin/sdk';

const client = new SetokinClient({
  apiKey: 'your-api-key',
  baseURL: 'https://api.setokin.com/v1'
});

// Create stock in
const stockIn = await client.stockIn.create({
  item_id: 'aa0e8400-e29b-41d4-a716-446655440005',
  quantity: 10.0,
  purchase_date: '2024-03-11',
  expiry_date: '2024-03-18'
});

// Get current inventory
const inventory = await client.inventory.list({
  low_stock_only: true
});
```

---

## Testing

### Sandbox Environment

Test API without affecting production data.

**Base URL:** `https://sandbox-api.setokin.com/v1`

### Test Credentials

```
Email: test@setokin.com
Password: TestPass123!
```

### Test Data

Sandbox is pre-populated with sample data:
- 50 items across 8 categories
- 100 batches with various expiry dates
- 200 stock transactions

### Reset Sandbox

```
POST /sandbox/reset
```

Resets sandbox to initial state.

---


## Best Practices

### Authentication

1. **Store tokens securely**
   - Access token in memory or secure storage
   - Refresh token in httpOnly cookie
   - Never expose tokens in URLs or logs

2. **Handle token expiration**
   - Implement automatic token refresh
   - Retry failed requests after refresh
   - Logout user if refresh fails

3. **Use HTTPS only**
   - Never send tokens over HTTP
   - Validate SSL certificates

### Error Handling

1. **Check HTTP status codes**
   - 2xx = Success
   - 4xx = Client error (fix request)
   - 5xx = Server error (retry with backoff)

2. **Parse error responses**
   - Always check `error.code` for specific handling
   - Display `error.message` to users
   - Log `error.details` for debugging

3. **Implement retry logic**
   - Retry 5xx errors with exponential backoff
   - Don't retry 4xx errors (except 429)
   - Maximum 3 retry attempts

### Performance

1. **Use pagination**
   - Always set reasonable `limit` values
   - Use cursor-based pagination for large datasets
   - Cache results when appropriate

2. **Request only needed fields**
   - Use `fields` parameter to reduce payload
   - Minimize nested object expansion

3. **Batch operations**
   - Group related operations when possible
   - Use bulk endpoints (when available)

### Data Validation

1. **Validate before sending**
   - Check required fields
   - Validate data types and formats
   - Verify business logic constraints

2. **Handle validation errors**
   - Parse `error.details` for field-specific errors
   - Display errors next to relevant form fields
   - Provide clear error messages to users

---

## Support & Resources

### Documentation

- **API Reference:** https://docs.setokin.com/api
- **Guides:** https://docs.setokin.com/guides
- **Changelog:** https://docs.setokin.com/changelog

### Community

- **Discord:** https://discord.gg/setokin
- **GitHub:** https://github.com/setokin
- **Forum:** https://community.setokin.com

### Support

- **Email:** support@setokin.com
- **Status Page:** https://status.setokin.com
- **Response Time:** < 24 hours

---

## Changelog

### v1.0.0 (2024-03-11)

**Initial Release**

- Authentication with JWT dual token
- Item management (CRUD)
- Stock in/out with FEFO logic
- Batch tracking and expiry alerts
- Inventory dashboard
- Daily/Weekly/Monthly reports
- File upload with presigned URLs
- Rate limiting
- Comprehensive error handling

---

## Appendix

### Date & Time Formats

All dates and times use ISO 8601 format:

- **Date:** `YYYY-MM-DD` (e.g., `2024-03-11`)
- **DateTime:** `YYYY-MM-DDTHH:mm:ssZ` (e.g., `2024-03-11T10:30:00Z`)
- **Timezone:** Always UTC (Z suffix)

### UUID Format

All IDs use UUID v4 format:

```
550e8400-e29b-41d4-a716-446655440000
```

### Decimal Precision

Quantity and price fields use decimal precision:

- **Quantity:** Up to 3 decimal places (e.g., `10.500`)
- **Price:** Up to 2 decimal places (e.g., `150000.00`)

### Allowed File Types

| Purpose | Allowed Types | Max Size |
|---------|---------------|----------|
| `stock_in_invoice` | PDF, JPEG, PNG | 10 MB |
| `item_image` | JPEG, PNG, WebP | 5 MB |
| `supplier_document` | PDF, DOC, DOCX | 10 MB |

---

**End of API Documentation**

For questions or feedback, contact: api@setokin.com
