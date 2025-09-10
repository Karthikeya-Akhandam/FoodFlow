# FoodFlow API Documentation

## Overview

FoodFlow is a food donation platform API that connects organizations offering surplus food with organizations that need it. The platform uses a credit-based system and matching algorithm to facilitate efficient food distribution.

**Base URL**: `http://localhost:8080`  
**API Version**: v1

## Authentication

The API uses JWT Bearer token authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-access-token>
```

### Token Expiration
- **Access Token**: 24 hours (configurable)
- **Refresh Token**: 720 hours (30 days, configurable)

## User Roles

- **ADMIN**: System administration and oversight
- **ORG**: Organizations that can create claims and manage credits
- **COLLAB**: Collaborators that can create donation offers and earn tokens

## Error Handling

All errors follow a consistent format:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message",
    "details": {}
  }
}
```

### Common Error Codes
- `VALIDATION_ERROR`: Request validation failed
- `UNAUTHORIZED`: Authentication required
- `FORBIDDEN`: Access denied
- `NOT_FOUND`: Resource not found
- `CONFLICT`: Resource conflict
- `INTERNAL_ERROR`: Server error
- `RATE_LIMITED`: Too many requests
- `IDEMPOTENCY_KEY_ERROR`: Idempotency key issues

## Pagination

List endpoints support pagination with these query parameters:

- `page`: Page number (default: 1)
- `limit`: Items per page (default: 20, max: 100)

Response format:
```json
{
  "items": [],
  "page": 1,
  "limit": 20,
  "total": 100
}
```

## Idempotency

Critical operations support idempotency using the `Idempotency-Key` header:

```
Idempotency-Key: <unique-key>
```

## Health Check

### GET /health

Returns system health status.

**Response:**
```json
{
  "status": "ok",
  "timestamp": "2024-01-01T00:00:00Z",
  "service": "foodflow-api"
}
```

---

## Authentication Endpoints

### POST /v1/auth/signup

Register a new user.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "role": "ORG"
}
```

**Validations:**
- `email`: Required, valid email format
- `password`: Required, minimum 8 characters
- `role`: Required, one of: `ORG`, `COLLAB`

**Response (201):**
```json
{
  "user_id": "uuid",
  "email": "user@example.com",
  "role": "ORG",
  "status": "ACTIVE"
}
```

**Errors:**
- 400: Validation error
- 409: Email already exists
- 500: Internal server error

### POST /v1/auth/login

Authenticate a user.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (200):**
```json
{
  "access_token": "jwt-token",
  "expires_in": 86400,
  "token_type": "Bearer"
}
```

**Errors:**
- 400: Validation error
- 401: Invalid credentials
- 500: Internal server error

### GET /v1/auth/me

Get current user information. **Requires authentication.**

**Response (200):**
```json
{
  "user_id": "uuid",
  "email": "user@example.com",
  "role": "ORG",
  "status": "ACTIVE"
}
```

**Errors:**
- 401: Unauthorized
- 404: User not found
- 500: Internal server error

---

## Profile Management

### GET /v1/profile

Get current user's profile. **Requires authentication.**

**Response (200):**
```json
{
  "name": "John Doe",
  "phone": "9876543210",
  "pincode": "123456",
  "city": "Mumbai",
  "district": "Mumbai",
  "state": "Maharashtra",
  "address": "123 Main St",
  "is_remote_org": false
}
```

### PUT /v1/profile

Update current user's profile. **Requires authentication.**

**Request Body (all fields optional):**
```json
{
  "name": "John Doe",
  "phone": "9876543210",
  "pincode": "123456",
  "city": "Mumbai",
  "district": "Mumbai",
  "state": "Maharashtra",
  "address": "123 Main St",
  "is_remote_org": false
}
```

**Validations:**
- `name`: 2-100 characters
- `phone`: Exactly 10 digits
- `pincode`: Exactly 6 characters
- `city`: 2-50 characters
- `district`: 2-50 characters
- `state`: 2-50 characters
- `address`: Maximum 200 characters

**Response (200):** Updated profile object

---

## Organization Management

### POST /v1/organizations

Create an organization profile. **Requires authentication (ORG role).**

**Request Body:**
```json
{
  "org_type": "NGO",
  "purpose_focus": ["CHILDREN", "ELDERLY"],
  "last_month_people_fed": 500
}
```

**Validations:**
- `org_type`: Required, 2-50 characters
- `purpose_focus`: Required array, values: `CHILDREN`, `ELDERLY`, `WOMEN`, `GENERAL`, `EMERGENCY`
- `last_month_people_fed`: Required, minimum 0

**Response (201):**
```json
{
  "org_id": "uuid",
  "org_type": "NGO",
  "purpose_focus": ["CHILDREN", "ELDERLY"],
  "last_month_people_fed": 500
}
```

### GET /v1/organizations/:id

Get organization details. **Requires authentication.**

**Response (200):** Organization object

### PUT /v1/organizations/:id

Update organization. **Requires authentication (Owner or Admin).**

**Request Body (all fields optional):**
```json
{
  "org_type": "NGO",
  "purpose_focus": ["CHILDREN"],
  "last_month_people_fed": 600
}
```

### GET /v1/organizations

List organizations. **Requires authentication.**

**Query Parameters:**
- `page`: Page number
- `limit`: Items per page

**Response (200):** Paginated list of organizations

---

## Collaborator Management

### POST /v1/collaborators

Create a collaborator profile. **Requires authentication (COLLAB role).**

**Request Body:**
```json
{
  "collab_type": "Individual Donor"
}
```

**Validations:**
- `collab_type`: Required, 2-50 characters

**Response (201):**
```json
{
  "collab_id": "uuid",
  "collab_type": "Individual Donor"
}
```

### GET /v1/collaborators/:id

Get collaborator details. **Requires authentication.**

### PUT /v1/collaborators/:id

Update collaborator. **Requires authentication (Owner or Admin).**

### GET /v1/collaborators

List collaborators. **Requires authentication.**

---

## Donation Offers

### POST /v1/offers

Create a donation offer. **Requires authentication (COLLAB role).**

**Request Body:**
```json
{
  "title": "Fresh vegetables from restaurant",
  "description": "Mixed vegetables, good for 50 people",
  "ready_from": "2024-01-01T10:00:00Z",
  "expires_at": "2024-01-01T18:00:00Z",
  "estimated_servings": 50,
  "purpose": "GENERAL",
  "pincode": "123456",
  "city": "Mumbai",
  "state": "Maharashtra"
}
```

**Validations:**
- `title`: Required, 5-100 characters
- `description`: Optional, max 500 characters
- `ready_from`: Required, ISO timestamp
- `expires_at`: Required, must be after `ready_from`
- `estimated_servings`: Required, 1-10000
- `purpose`: Optional, one of: `CHILDREN`, `ELDERLY`, `WOMEN`, `GENERAL`, `EMERGENCY`
- `pincode`: Optional, exactly 6 characters
- `city`: Optional, 2-50 characters
- `state`: Optional, 2-50 characters

**Response (201):**
```json
{
  "offer_id": "uuid",
  "status": "OPEN"
}
```

### GET /v1/offers/:id

Get offer details. **Requires authentication.**

**Response (200):**
```json
{
  "offer_id": "uuid",
  "collab_id": "uuid",
  "title": "Fresh vegetables",
  "description": "Mixed vegetables",
  "ready_from": "2024-01-01T10:00:00Z",
  "expires_at": "2024-01-01T18:00:00Z",
  "estimated_servings": 50,
  "purpose": "GENERAL",
  "location": {
    "pincode": "123456",
    "city": "Mumbai",
    "state": "Maharashtra"
  },
  "status": "OPEN"
}
```

### PUT /v1/offers/:id

Update an offer. **Requires authentication (Owner only).**

**Request Body (all fields optional):**
```json
{
  "title": "Updated title",
  "description": "Updated description",
  "ready_from": "2024-01-01T11:00:00Z",
  "expires_at": "2024-01-01T19:00:00Z",
  "estimated_servings": 60,
  "purpose": "CHILDREN",
  "pincode": "123457",
  "city": "Delhi",
  "state": "Delhi"
}
```

### GET /v1/offers

List offers. **Requires authentication.**

**Query Parameters:**
- `status`: Filter by status (`OPEN`, `PENDING_CONFIRM`, `CLAIMED`, `EXPIRED`, `CANCELLED`)
- `page`: Page number
- `limit`: Items per page

**Access Control:**
- COLLABs see only their own offers
- ADMINs see all offers

**Response (200):** Paginated list of offers

### GET /v1/offers/nearby

Get nearby offers with priority scoring. **Requires authentication (ORG role).**

**Query Parameters:**
- `purpose`: Filter by purpose
- `page`: Page number
- `limit`: Items per page

**Response (200):**
```json
{
  "items": [
    {
      "offer_id": "uuid",
      "title": "Fresh vegetables",
      "estimated_servings": 50,
      "purpose": "GENERAL",
      "ready_from": "2024-01-01T10:00:00Z",
      "expires_at": "2024-01-01T18:00:00Z",
      "score": 85.5,
      "proximity_tier": "pincode"
    }
  ],
  "page": 1,
  "limit": 20,
  "total": 5
}
```

---

## Claims Management

### POST /v1/claims

Create a claim on an offer. **Requires authentication (ORG role).**

**Request Body:**
```json
{
  "requested_servings": 30,
  "on_behalf_remote_org_id": "uuid"
}
```

**Validations:**
- `requested_servings`: Required, 1-10000
- `on_behalf_remote_org_id`: Optional UUID for remote org claims

**Response (201):**
```json
{
  "claim_id": "uuid",
  "offer_id": "uuid",
  "status": "REQUESTED",
  "priority_score": 75.5,
  "requested_servings": 30,
  "created_at": "2024-01-01T12:00:00Z"
}
```

### GET /v1/claims/:id

Get claim details. **Requires authentication.**

### PUT /v1/claims/:id/status

Update claim status. **Requires authentication (Owner or Admin).**

**Request Body:**
```json
{
  "status": "WON"
}
```

**Valid statuses:** `REQUESTED`, `WON`, `LOST`, `CANCELLED`

### GET /v1/claims

List claims. **Requires authentication.**

**Query Parameters:**
- `status`: Filter by status
- `page`: Page number
- `limit`: Items per page

---

## Redemptions

### POST /v1/redemptions

Confirm a redemption. **Requires authentication (ORG role).**

**Request Body:**
```json
{
  "offer_id": "uuid",
  "servings_accepted": 25
}
```

**Validations:**
- `offer_id`: Required UUID
- `servings_accepted`: Required, 1-10000

**Response (201):**
```json
{
  "redemption_id": "uuid",
  "credits_spent": 25,
  "tokens_awarded": 5,
  "offer_status": "CLAIMED"
}
```

### GET /v1/redemptions/:id

Get redemption details. **Requires authentication.**

### PUT /v1/redemptions/:id/status

Update redemption status. **Requires authentication (Owner or Admin).**

**Request Body:**
```json
{
  "status": "CONFIRMED"
}
```

**Valid statuses:** `PENDING`, `CONFIRMED`, `CANCELLED`

### GET /v1/redemptions

List redemptions. **Requires authentication.**

---

## Credits System

### GET /v1/credits

Get organization credits. **Requires authentication (ORG role).**

**Response (200):**
```json
{
  "total_credits": 150,
  "monthly_credits": [
    {
      "month": "2024-01",
      "credits_issued": 100,
      "credits_remaining": 75,
      "expires_at": "2024-01-31T23:59:59Z"
    }
  ]
}
```

### GET /v1/credits/history

Get credit transaction history. **Requires authentication (ORG role).**

**Response (200):**
```json
{
  "items": [
    {
      "transaction_id": "uuid",
      "amount": 25,
      "type": "SPEND",
      "description": "Claimed food donation",
      "created_at": "2024-01-01T12:00:00Z"
    }
  ],
  "page": 1,
  "limit": 20,
  "total": 10
}
```

### POST /v1/credits/spend

Spend credits manually. **Requires authentication (ORG role).**

**Request Body:**
```json
{
  "amount": 10,
  "reason": "Manual adjustment for special program"
}
```

**Validations:**
- `amount`: Required, minimum 1
- `reason`: Required, 5-200 characters

---

## Tokens System

### GET /v1/tokens

Get collaborator tokens. **Requires authentication (COLLAB role).**

**Response (200):**
```json
{
  "total_tokens": 50,
  "monthly_tokens": [
    {
      "month": "2024-01",
      "tokens_earned": 20
    }
  ]
}
```

### GET /v1/tokens/history

Get token transaction history. **Requires authentication (COLLAB role).**

**Response (200):**
```json
{
  "items": [
    {
      "transaction_id": "uuid",
      "amount": 5,
      "type": "EARN",
      "description": "Offer redeemed successfully",
      "created_at": "2024-01-01T12:00:00Z"
    }
  ],
  "page": 1,
  "limit": 20,
  "total": 8
}
```

### POST /v1/tokens/redeem

Redeem tokens for credits or other rewards. **Requires authentication (COLLAB role).**

**Request Body:**
```json
{
  "amount": 10,
  "reason": "Convert to credits for partner organization"
}
```

**Response (200):**
```json
{
  "redemption_id": "uuid",
  "tokens_spent": 10,
  "credits_earned": 2,
  "status": "CONFIRMED"
}
```

---

## Remote Organizations

### POST /v1/remote-orgs

Create a remote organization. **Requires authentication (ORG role).**

**Request Body:**
```json
{
  "name": "Rural Community Center",
  "contact_email": "contact@center.org",
  "contact_phone": "9876543210",
  "address": "123 Village Road, Block XYZ",
  "pincode": "123456",
  "city": "Rural Town",
  "state": "State Name",
  "serving_capacity": 200,
  "purpose_focus": ["CHILDREN", "WOMEN"]
}
```

**Validations:**
- `name`: Required, 2-100 characters
- `contact_email`: Required, valid email
- `contact_phone`: Optional, exactly 10 digits
- `address`: Required, 10-200 characters
- `pincode`: Required, exactly 6 characters
- `city`: Required, 2-50 characters
- `state`: Required, 2-50 characters
- `serving_capacity`: Required, minimum 1
- `purpose_focus`: Required array of purposes

### GET /v1/remote-orgs/:id

Get remote organization details. **Requires authentication.**

### PUT /v1/remote-orgs/:id

Update remote organization. **Requires authentication (Owner or Admin).**

### GET /v1/remote-orgs

List remote organizations. **Requires authentication.**

### POST /v1/remote-orgs/assign

Assign remote organization to proxy. **Requires authentication (ORG role).**

**Request Body:**
```json
{
  "remote_org_id": "uuid",
  "notes": "Assigned for better coverage in rural area"
}
```

---

## Admin Endpoints

All admin endpoints require **ADMIN role**.

### GET /v1/admin/stats

Get system statistics.

**Response (200):**
```json
{
  "total_users": 1250,
  "active_users": 1200,
  "total_organizations": 800,
  "total_collaborators": 400,
  "total_offers": 5000,
  "active_offers": 150,
  "total_claims": 3000,
  "total_redemptions": 2500
}
```

### GET /v1/admin/users/stats

Get user statistics.

**Response (200):**
```json
{
  "new_users_today": 5,
  "new_users_this_week": 35,
  "new_users_this_month": 120,
  "active_users_today": 150,
  "suspended_users": 10
}
```

### GET /v1/admin/donations/stats

Get donation statistics.

**Response (200):**
```json
{
  "offers_today": 15,
  "offers_this_week": 85,
  "offers_this_month": 350,
  "servings_offered": 25000,
  "servings_redeemed": 20000,
  "redemption_rate": 0.8,
  "avg_servings_per_offer": 71.4
}
```

### GET /v1/admin/users

List users with filters.

**Query Parameters:**
- `role`: Filter by role
- `status`: Filter by status
- `page`: Page number
- `limit`: Items per page

### PUT /v1/admin/users/:id/status

Update user status.

**Request Body:**
```json
{
  "status": "SUSPENDED"
}
```

**Valid statuses:** `ACTIVE`, `SUSPENDED`

### GET /v1/admin/audit-logs

Get audit logs.

**Query Parameters:**
- `action`: Filter by action
- `entity_type`: Filter by entity type
- `user_id`: Filter by user ID
- `page`: Page number
- `limit`: Items per page

**Response (200):**
```json
{
  "items": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "action": "CREATE",
      "entity_type": "OFFER",
      "entity_id": "uuid",
      "changes": {
        "title": "New offer created",
        "servings": 50
      },
      "ip_address": "192.168.1.1",
      "user_agent": "Mozilla/5.0...",
      "created_at": "2024-01-01T12:00:00Z"
    }
  ],
  "page": 1,
  "limit": 20,
  "total": 500
}
```

---

## Rate Limiting

The API implements rate limiting to prevent abuse:

- **Default**: 100 requests per second with burst of 200
- **Rate limit headers** are included in responses:
  - `X-RateLimit-Limit`: Requests allowed per window
  - `X-RateLimit-Remaining`: Remaining requests in window
  - `X-RateLimit-Reset`: Time when window resets

When rate limit is exceeded, you'll receive a `429 Too Many Requests` error.

---

## Data Types and Enums

### User Roles
- `ADMIN`: System administrator
- `ORG`: Organization
- `COLLAB`: Collaborator

### User Status
- `ACTIVE`: Active user
- `SUSPENDED`: Suspended user

### Donation Status
- `OPEN`: Available for claims
- `PENDING_CONFIRM`: Awaiting confirmation
- `CLAIMED`: Successfully claimed
- `EXPIRED`: Past expiration time
- `CANCELLED`: Cancelled by owner

### Claim Status
- `REQUESTED`: Claim submitted
- `WON`: Claim won through matching
- `LOST`: Claim lost to another organization
- `CANCELLED`: Claim cancelled

### Purpose Types
- `CHILDREN`: Food for children
- `ELDERLY`: Food for elderly
- `WOMEN`: Food for women
- `GENERAL`: General purpose
- `EMERGENCY`: Emergency distribution

### Redemption Status
- `PENDING`: Awaiting confirmation
- `CONFIRMED`: Confirmed and completed
- `CANCELLED`: Cancelled redemption

### Transaction Types

**Credit Transactions:**
- `ISSUE`: Credits issued
- `SPEND`: Credits spent
- `REFUND`: Credits refunded

**Token Transactions:**
- `EARN`: Tokens earned
- `REDEEM`: Tokens redeemed

---

## Matching Algorithm

The platform uses a sophisticated matching algorithm that considers:

### Priority Factors
1. **Geographic Proximity** (40% weight)
   - Same pincode: Highest priority
   - Same city: Medium priority
   - Same state: Lower priority

2. **Purpose Alignment** (25% weight)
   - Exact purpose match: Bonus points
   - Complementary purposes: Partial points

3. **Credit Pressure** (20% weight)
   - Organizations with fewer credits get higher priority
   - Prevents credit hoarding

4. **Reliability Score** (10% weight)
   - Based on successful redemption history
   - Consistent organizations get priority

5. **Remote Organization Bonus** (5% weight)
   - Remote organizations get extra points
   - Encourages service to underserved areas

### Matching Process
1. Organizations create claims for offers
2. System calculates priority scores for each claim
3. Highest scoring claim wins when offer is assigned
4. Credits are held during pending confirmation
5. Successful redemption transfers credits and awards tokens

---

## Best Practices

### Authentication
- Store tokens securely (use secure storage)
- Implement token refresh logic
- Handle 401 responses by redirecting to login

### Error Handling
- Always check response status codes
- Parse error responses for detailed messages
- Implement retry logic for 5xx errors
- Respect rate limiting (watch for 429 responses)

### Idempotency
- Use idempotency keys for critical operations
- Generate unique keys (UUID recommended)
- Store keys to prevent duplicate requests

### Pagination
- Always handle pagination for list endpoints
- Use reasonable page sizes (20-50 items)
- Cache results where appropriate

### Performance
- Use filters to reduce response sizes
- Implement client-side caching for reference data
- Batch related requests where possible

---

## Common Integration Patterns

### Organization Workflow
1. Register user with `ORG` role
2. Complete profile with location information
3. Create organization profile
4. Browse nearby offers
5. Create claims for suitable offers
6. Confirm redemptions when food is collected
7. Monitor credit balance and history

### Collaborator Workflow
1. Register user with `COLLAB` role
2. Complete profile with location information
3. Create collaborator profile
4. Create donation offers when food is available
5. Monitor offer status and claims
6. Confirm redemptions when food is collected
7. Track earned tokens and redeem when needed

### Admin Workflow
1. Monitor system statistics
2. Review user activity and audit logs
3. Manage user statuses (suspend/activate)
4. Oversee donation statistics and trends
5. Handle escalated issues and disputes

---

## WebSocket Integration (Future)

The API is designed to support real-time features through WebSocket connections:

### Planned Events
- Offer status changes
- New claims on offers
- Redemption confirmations
- Credit/token balance updates
- System notifications

### Connection Pattern
```javascript
const ws = new WebSocket('ws://localhost:8080/ws');
ws.onopen = () => {
  ws.send(JSON.stringify({
    type: 'auth',
    token: 'your-jwt-token'
  }));
};
```

---

This documentation provides comprehensive coverage of the FoodFlow API. For additional support or questions about implementation, please refer to the source code or contact the development team.