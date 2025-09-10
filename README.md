# FoodFlow - Food Donation Platform

A production-ready backend API for managing food donations between collaborators (restaurants, hotels, etc.) and organizations (NGOs, orphanages, shelters).

## Features

- **User Management**: Role-based authentication (Admin, Organization, Collaborator)
- **Donation Offers**: Collaborators can post food donation offers
- **Smart Matching**: AI-powered matching system based on proximity, purpose, and credits
- **Credit System**: Monthly credit allocation for organizations
- **Token Rewards**: Token system for collaborators based on donations
- **Remote Organizations**: Support for remote organizations with proxy assignments
- **Background Jobs**: Automated credit issuance, offer cleanup, and auto-assignment
- **Rate Limiting**: IP and user-based rate limiting
- **Idempotency**: Safe retry mechanisms for critical operations
- **Audit Trail**: Complete event logging for compliance

## Architecture

### Tech Stack
- **Backend**: Go 1.21 with Gin framework
- **Database**: PostgreSQL with sqlc for type-safe queries
- **Cache**: Redis for rate limiting, idempotency, and caching
- **Authentication**: JWT with RS256 signing
- **Background Jobs**: Cron-based job scheduling
- **Containerization**: Docker with docker-compose

### Project Structure
```
├── cmd/api/                 # Application entry point
├── config/                  # Configuration management
├── internal/
│   ├── core/               # Domain entities, DTOs, and business logic
│   │   ├── entities.go     # Core domain models
│   │   ├── dto.go         # Request/Response DTOs
│   │   ├── errors.go      # Error handling
│   │   ├── consts.go      # Constants and enums
│   │   ├── services/      # Business logic services
│   │   └── repos/         # Repository interfaces
│   ├── db/                # Database layer
│   │   ├── migrations/    # SQL migrations
│   │   └── sqlc/         # Generated SQL code
│   ├── http/              # HTTP layer
│   │   ├── handlers/      # HTTP handlers
│   │   ├── middleware/    # Middleware (auth, rate limiting, etc.)
│   │   └── router.go     # Route configuration
│   ├── jobs/              # Background jobs
│   └── lib/               # Shared utilities
├── docker-compose.yml     # Development environment
├── Dockerfile            # Production container
└── Makefile             # Development commands
```

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Go 1.21+ (for local development)
- Make (optional, for convenience commands)

### Development Setup

1. **Clone and setup**:
   ```bash
   git clone <repository-url>
   cd FoodFlow
   make setup
   ```

2. **Start services**:
   ```bash
   make docker-compose-up
   ```

3. **Start the API server** (auto-migration runs automatically):
   ```bash
   make run
   ```

The API will be available at `http://localhost:8080`

### Production Deployment

1. **Build and run with Docker**:
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

2. **Or build the binary**:
   ```bash
   make build
   ./bin/foodflow
   ```

## API Documentation

### Authentication

All protected endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <jwt-token>
```

### Core Endpoints

#### Authentication
- `POST /v1/auth/signup` - Register new user
- `POST /v1/auth/login` - Login user
- `GET /v1/auth/me` - Get current user info

#### Offers (Collaborators)
- `POST /v1/offers` - Create donation offer
- `GET /v1/offers` - List offers (own offers for collaborators)
- `GET /v1/offers/:id` - Get offer details
- `POST /v1/offers/:id/cancel` - Cancel offer

#### Offers (Organizations)
- `GET /v1/offers/nearby` - Get nearby offers with priority scores

#### Claims (Organizations)
- `POST /v1/offers/:id/claim` - Claim an offer
- `GET /v1/claims` - List organization's claims
- `POST /v1/claims/:id/cancel` - Cancel claim

#### Redemption
- `POST /v1/offers/:id/redeem` - Complete redemption (spend credits, award tokens)

#### Credits & Tokens
- `GET /v1/org/credits/current` - Get current month credits
- `GET /v1/collab/tokens/current` - Get current month tokens

### Request/Response Examples

#### Create Offer
```bash
curl -X POST http://localhost:8080/v1/offers \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Veg Biryani Trays",
    "description": "Fresh vegetarian biryani, no onion/garlic",
    "ready_from": "2025-01-15T12:30:00Z",
    "expires_at": "2025-01-15T16:30:00Z",
    "estimated_servings": 120,
    "purpose": "CHILDREN",
    "pincode": "560001",
    "city": "Bengaluru",
    "state": "Karnataka"
  }'
```

#### Claim Offer
```bash
curl -X POST http://localhost:8080/v1/offers/{offer-id}/claim \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "requested_servings": 80
  }'
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DB_URL` | PostgreSQL connection string | `postgres://foodflow:password@localhost:5432/foodflow?sslmode=disable` |
| `REDIS_URL` | Redis connection string | `redis://localhost:6379` |
| `JWT_PRIVATE_KEY_PATH` | Path to JWT private key | `./keys/app.rsa` |
| `JWT_PUBLIC_KEY_PATH` | Path to JWT public key | `./keys/app.rsa.pub` |
| `SERVINGS_PER_CREDIT` | Servings per credit | `10` |
| `TOKEN_MULTIPLIER` | Token multiplier for rewards | `5` |
| `CLAIM_COOLOFF_SECONDS` | Cool-off period for auto-assignment | `600` |
| `RATE_LIMIT_RPS` | Rate limit requests per second | `10` |
| `RATE_LIMIT_BURST` | Rate limit burst capacity | `20` |

### JWT Keys

Generate JWT keys for development:
```bash
make keys
```

## Business Logic

### Matching Algorithm

The system uses a sophisticated scoring algorithm to match offers with organizations:

```
Score = 40 * proximity_weight + 25 * purpose_match_weight + 
        20 * credit_pressure_weight + 10 * reliability_weight + 
        5 * remote_proxy_bonus
```

- **Proximity**: Same pincode (1.0), same city (0.7), same state (0.4)
- **Purpose Match**: 1.0 if organization's purpose focus matches offer purpose
- **Credit Pressure**: Higher weight for organizations with fewer credits
- **Reliability**: Based on historical performance
- **Remote Proxy**: Bonus for organizations serving remote areas

### Credit System Details

Organizations receive monthly credits based on their capacity:
```
Monthly Credits = BASE_CREDITS + (K_FACTOR × last_month_people_fed) + remote_bonus
```

Configuration:
- `MONTHLY_CREDITS_BASE=50` - Base credits for all organizations
- `MONTHLY_CREDITS_K=0.2` - Multiplier based on people served
- `REMOTE_PROXY_BONUS=0.1` - Additional 10% for remote organizations
- `SERVINGS_PER_CREDIT=10` - How many servings each credit covers

**Example**: An organization that fed 200 people gets:
- Base: 50 credits
- Performance: 200 × 0.2 = 40 credits  
- Total: 90 credits (99 if remote)
- Can claim: 900 servings (990 if remote)

**Credit Rules**:
- Credits expire at the end of each month (no carry-over in v1)
- 1 credit = 10 servings (configurable)

### Token System

- Collaborators earn tokens when their offers are redeemed
- Tokens = `credits_spent * TOKEN_MULTIPLIER`
- Tokens are tracked monthly

## Background Jobs

### Monthly Credits Job
- Runs on the 1st of every month at 00:05 UTC
- Issues new credits to all organizations
- Expires previous month's credits

### Offer Cleanup Job
- Runs every hour
- Expires offers past their expiration time
- Cancels stale claims (>24 hours old)

### Auto-Assignment Job
- Runs every 2 minutes
- Automatically assigns offers to highest-scoring claims after cool-off period
- Only for offers with purpose specified

## Development

### Available Commands

```bash
make help                 # Show all available commands
make build               # Build the application
make run                 # Run locally
make test                # Run tests
make docker-compose-up   # Start all services
make migrate-up          # Run database migrations
make seed                # Seed sample data
make sqlc-generate       # Generate SQL code
```

### Database Migrations

**Auto-migration**: Migrations run automatically when the application starts.

**Manual migration commands** (for development):
```bash
make migrate-up          # Apply migrations manually
make migrate-down        # Rollback migrations manually
make migrate-force VERSION=1  # Force migration version
```

### Testing

```bash
make test                # Run all tests
make test-coverage       # Run tests with coverage
```

## Monitoring

### Health Check
```bash
curl http://localhost:8080/health
```

### Logs
```bash
make logs                # View application logs
docker-compose logs -f   # View all service logs
```

## Security Features

- **JWT Authentication**: RS256 signed tokens
- **Password Hashing**: Argon2id for secure password storage
- **Rate Limiting**: IP and user-based rate limiting
- **Idempotency**: Safe retry mechanisms
- **Input Validation**: Comprehensive request validation
- **SQL Injection Protection**: sqlc generated queries
- **CORS**: Configurable cross-origin resource sharing

## Performance

- **Connection Pooling**: Optimized database connection management
- **Redis Caching**: Fast access to frequently used data
- **Background Jobs**: Non-blocking operations
- **Efficient Queries**: Optimized SQL with proper indexing
- **Graceful Shutdown**: Clean application termination

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support and questions, please open an issue in the repository or contact the development team.
