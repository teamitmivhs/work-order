# API Reference - Frontend Integration Guide

## Base URL
```
http://localhost:8080/api
```

---

## Authentication Flow

### 1. Register
Create a new user account.

**Endpoint**: `POST /user/register`

**Request**:
```json
{
  "name": "username",
  "password": "Password123"
}
```

**Password Requirements**:
- Minimum 8 characters
- Must contain uppercase letter (A-Z)
- Must contain lowercase letter (a-z)
- Must contain digit (0-9)

**Response** (201 Created):
```json
{
  "message": "Registration successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "member": {
      "id": 1,
      "name": "username",
      "role": "Operator",
      "status": "standby"
    }
  }
}
```

**Error Responses**:
- 400 Bad Request: Invalid input or weak password
- 409 Conflict: Username already exists

---

### 2. Login
Authenticate and get JWT token.

**Endpoint**: `POST /user/login`

**Request**:
```json
{
  "name": "username",
  "password": "Password123"
}
```

**Response** (200 OK):
```json
{
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "member": {
      "id": 1,
      "name": "username",
      "role": "Operator",
      "status": "standby"
    }
  }
}
```

**Error Responses**:
- 401 Unauthorized: Invalid credentials

---

### 3. Get Profile
Retrieve current user profile.

**Endpoint**: `GET /user/profile`

**Headers**:
```
Authorization: Bearer <token>
```

**Response** (200 OK):
```json
{
  "id": 1,
  "name": "username",
  "role": "Operator",
  "status": "standby",
  "avatar": ""
}
```

**Error Responses**:
- 401 Unauthorized: Missing or invalid token
- 404 Not Found: User not found

---

## Members Management

### Get All Members
Retrieve list of all team members (public endpoint).

**Endpoint**: `GET /members`

**Response** (200 OK):
```json
[
  {
    "id": 1,
    "name": "john_doe",
    "role": "Operator",
    "status": "standby",
    "avatar": "avatar_url"
  },
  {
    "id": 2,
    "name": "jane_doe",
    "role": "Admin",
    "status": "onjob",
    "avatar": "avatar_url"
  }
]
```

---

## Work Order Management

### 1. Create Work Order
Create a new work order.

**Endpoint**: `POST /workorders`

**Headers**:
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request**:
```json
{
  "priority": "high",
  "requester": "Admin",
  "location": "Building A - Room 101",
  "device": "Printer",
  "problem": "Paper jam in tray 2",
  "status": "pending",
  "time_display": "2024-12-27 10:30"
}
```

**Priority Values**: `low`, `medium`, `high`, `urgent`

**Response** (201 Created):
```json
{
  "message": "Work order created successfully",
  "data": {
    "id": 123
  }
}
```

**Error Responses**:
- 400 Bad Request: Missing required fields or invalid priority
- 401 Unauthorized: Missing or invalid token
- 500 Internal Server Error: Database error

---

### 2. Get Work Orders
Retrieve work orders (filtered by user role).

**Endpoint**: `GET /workorders`

**Headers**:
```
Authorization: Bearer <token>
```

**Query Parameters** (optional):
```
?status=pending
?priority=high
```

**Response** (200 OK):
```json
[
  {
    "id": 1,
    "priority": "high",
    "time": "2024-12-27 10:30",
    "requester": "Admin",
    "location": "Building A",
    "device": "Printer",
    "problem": "Paper jam",
    "executors": [1, 2],
    "workingHours": 2,
    "status": "completed",
    "safetyChecklist": ["Check power", "Test device"],
    "completedAt": "2024-12-27 12:30"
  }
]
```

**Filtering Rules**:
- **Admin**: Sees all work orders
- **Operator**: Sees only assigned work orders (in `executors` list)

---

### 3. Take Work Order
Assign operators to work order and start progress.

**Endpoint**: `POST /workorders/{id}/take`

**Headers**:
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request**:
```json
{
  "status": "progress",
  "executors": [1, 2],
  "safety_checklist_items": ["Check power", "Check safety gear"]
}
```

**Response** (200 OK):
```json
{
  "message": "Order taken successfully",
  "data": {
    "id": 123
  }
}
```

**Error Responses**:
- 400 Bad Request: Status must be 'progress'
- 401 Unauthorized: Missing or invalid token
- 403 Forbidden: User not assigned to this work order
- 404 Not Found: Work order not found

---

### 4. Get Safety Checklist
Retrieve safety checklist for a work order.

**Endpoint**: `GET /workorders/{id}/checklist`

**Headers**:
```
Authorization: Bearer <token>
```

**Response** (200 OK):
```json
{
  "checklist": [
    "Check device power",
    "Verify functionality",
    "Test after repair",
    "Document actions taken"
  ]
}
```

---

### 5. Update Safety Checklist
Update safety checklist items for work order.

**Endpoint**: `PUT /workorders/{id}/checklist`

**Headers**:
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request**:
```json
{
  "checklist_items": [
    "Power checked - OK",
    "Functionality verified",
    "Repair completed",
    "Actions documented"
  ]
}
```

**Response** (200 OK):
```json
{
  "message": "Safety checklist updated successfully"
}
```

**Error Responses**:
- 400 Bad Request: Checklist items empty
- 401 Unauthorized: Missing or invalid token

---

### 6. Complete Work Order
Mark work order as completed.

**Endpoint**: `PATCH /workorders/{id}/complete`

**Headers**:
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request**:
```json
{
  "status": "completed",
  "completed_at_display": "2024-12-27 15:30"
}
```

**Response** (200 OK):
```json
{
  "message": "Order completed successfully",
  "data": {
    "id": 123
  }
}
```

**Validation**:
- Only assigned members can complete
- Safety checklist must be fulfilled
- Status must be 'completed'

**Error Responses**:
- 400 Bad Request: Safety checklist not completed or invalid status
- 401 Unauthorized: Missing or invalid token
- 403 Forbidden: User not assigned to this work order

---

### 7. Delete Work Order
Delete a work order (Admin only).

**Endpoint**: `DELETE /workorders/{id}`

**Headers**:
```
Authorization: Bearer <token>
```

**Response** (200 OK):
```json
{
  "message": "Order deleted successfully",
  "data": {
    "id": 123
  }
}
```

**Error Responses**:
- 401 Unauthorized: Missing or invalid token
- 403 Forbidden: Admin access required

---

## Performance & Analytics

### Get Kaizen Metrics
Retrieve performance and efficiency metrics.

**Endpoint**: `GET /kaizen`

**Headers**:
```
Authorization: Bearer <token>
```

**Response** (200 OK):
```json
{
  "totalKaizens": 50,
  "implementedKaizens": 35,
  "pendingKaizens": 15
}
```

---

## Error Response Format

All error responses follow this format:

```json
{
  "code": 400,
  "message": "Error message describing what went wrong",
  "details": "Optional additional context"
}
```

**Common HTTP Status Codes**:
- `200 OK`: Successful request
- `201 Created`: Resource successfully created
- `400 Bad Request`: Invalid input or validation failed
- `401 Unauthorized`: Missing or invalid authentication token
- `403 Forbidden`: User doesn't have permission for this action
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource already exists (e.g., duplicate username)
- `500 Internal Server Error`: Server-side error

---

## Authentication Header

All protected endpoints require the `Authorization` header with JWT token:

```
Authorization: Bearer <token>
```

**Token Storage** (Frontend):
```javascript
// Save token after login
localStorage.setItem('auth_token', response.data.token);

// Retrieve token for requests
const token = localStorage.getItem('auth_token');

// Use in fetch
fetch('/api/workorders', {
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
});
```

**Token Expiration**: 24 hours. After expiration, user needs to login again.

---

## Frontend Implementation Examples

### JavaScript - Fetch API

```javascript
// 1. Register
async function register(username, password) {
  const response = await fetch('/api/user/register', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name: username, password })
  });
  const data = await response.json();
  if (response.ok) {
    localStorage.setItem('auth_token', data.data.token);
    return data.data.member;
  }
  throw new Error(data.message);
}

// 2. Login
async function login(username, password) {
  const response = await fetch('/api/user/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name: username, password })
  });
  const data = await response.json();
  if (response.ok) {
    localStorage.setItem('auth_token', data.data.token);
    return data.data.member;
  }
  throw new Error(data.message);
}

// 3. Protected request
async function getWorkOrders() {
  const token = localStorage.getItem('auth_token');
  const response = await fetch('/api/workorders', {
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  });
  const data = await response.json();
  if (response.ok) {
    return data;
  }
  throw new Error(data.message);
}

// 4. Create work order
async function createWorkOrder(workOrder) {
  const token = localStorage.getItem('auth_token');
  const response = await fetch('/api/workorders', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(workOrder)
  });
  const data = await response.json();
  if (response.ok) {
    return data.data;
  }
  throw new Error(data.message);
}
```

---

## CORS Configuration

The backend is configured with CORS enabled:
- **Allowed Origins**: * (all)
- **Allowed Methods**: GET, POST, PUT, PATCH, DELETE, OPTIONS
- **Allowed Headers**: Origin, Content-Type, Authorization

---

## Rate Limiting

Currently not implemented. Recommended for production:
- 100 requests/minute per IP for public endpoints
- 50 requests/minute per user for authenticated endpoints

---

## Data Validation Rules

### Username
- Required
- 3-50 characters
- Alphanumeric recommended

### Password
- Required
- Minimum 8 characters
- Must include uppercase letter
- Must include lowercase letter
- Must include digit

### Priority
- Must be one of: `low`, `medium`, `high`, `urgent`

### Status
- Work orders: `pending`, `progress`, `completed`
- Members: `standby`, `onjob`

### Required Fields for Work Order
- `priority`
- `requester`
- `location`
- `device`
- `problem`

