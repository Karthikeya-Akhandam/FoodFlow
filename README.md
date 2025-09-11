# FoodFlow - Food Donation Platform

A full-stack application for managing food donations between collaborators (restaurants, hotels, etc.) and organizations (NGOs, orphanages, shelters).

## Project Structure

This is a monorepo containing both backend and frontend applications:

```
FoodFlow/
├── foodflow-backend/     # Go backend API (✅ Complete)
│   ├── cmd/             # Application entry point
│   ├── internal/        # Backend business logic
│   ├── config/          # Configuration management
│   ├── migrations/      # Database migrations
│   ├── go.mod          # Go module file
│   ├── .env            # Environment configuration
│   └── API_DOCUMENTATION.md  # Complete API docs
├── foodflow-frontend/   # React/Next.js frontend (🚧 Coming Soon)
│   └── (to be created)
├── README.md           # This file - project overview
├── CLAUDE.md          # Development context and architecture
├── LICENSE            # MIT License
└── Makefile          # Common development commands
```

## 🎉 Backend Status: Production Ready ✅

The backend API is **complete** with 50+ endpoints, comprehensive authentication, credit system, and full documentation ready for frontend integration.

## Features

- **User Management**: Role-based authentication (Admin, Organization, Collaborator)
- **Donation Offers**: Collaborators can post food donation offers
- **Smart Matching**: Algorithm-based matching system based on proximity, purpose, and credits
- **Credit System**: Monthly credit allocation for organizations
- **Token Rewards**: Token system for collaborators based on donations
- **Remote Organizations**: Support for remote organizations with proxy assignments
- **Admin Dashboard**: Complete admin functionality with system stats and user management
- **Rate Limiting**: IP and user-based rate limiting
- **Idempotency**: Safe retry mechanisms for critical operations
- **Audit Trail**: Complete event logging for compliance

## Backend Architecture

### Tech Stack
- **Backend**: Go 1.21 with Gin framework
- **Database**: PostgreSQL with type-safe queries
- **Cache**: Redis for rate limiting, idempotency, and caching
- **Authentication**: JWT with HMAC-SHA256 signing
- **Containerization**: Docker with docker-compose

### Backend Structure
```
foodflow-backend/
├── cmd/api/                 # Application entry point
├── config/                  # Configuration management
├── internal/
│   ├── core/               # Domain entities, DTOs, and business logic
│   │   ├── entities.go     # Core domain models
│   │   ├── dto.go         # Request/Response DTOs
│   │   ├── errors.go      # Error handling
│   │   ├── services/      # Business logic services
│   │   └── repos/         # Repository interfaces
│   ├── http/              # HTTP layer
│   │   ├── handlers/      # HTTP handlers
│   │   ├── middleware/    # Middleware (auth, rate limiting, etc.)
│   │   └── router.go     # Route configuration
│   ├── lib/               # Shared utilities
│   └── db/               # Database layer
├── migrations/            # SQL migrations
├── docker-compose.yml    # Development environment
├── Dockerfile           # Production container
└── API_DOCUMENTATION.md # Complete API documentation
```

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Go 1.21+ (for local development)

### Backend Development Setup

1. **Clone and setup**:
   ```bash
   git clone <repository-url>
   cd FoodFlow/foodflow-backend
   cp env.sample .env
   # Edit .env file with your database credentials
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Start services** (from project root):
   ```bash
   cd ..
   make docker-compose-up
   ```

4. **Start the API server** (auto-loads .env and runs migrations):
   ```bash
   cd foodflow-backend
   go run cmd/api/main.go
   ```

The API will be available at `http://localhost:8080`

### Backend Production Deployment

1. **Build and run with Docker**:
   ```bash
   cd foodflow-backend
   docker-compose -f docker-compose.prod.yml up -d
   ```

2. **Or build the binary**:
   ```bash
   cd foodflow-backend
   go build -o bin/foodflow cmd/api/main.go
   ./bin/foodflow
   ```

## API Documentation

📖 **Complete API documentation available**: See `foodflow-backend/API_DOCUMENTATION.md` for comprehensive endpoint documentation with examples.

### Authentication

All protected endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <jwt-token>
```

### Key API Features
- **50+ API endpoints** covering all functionality
- **Role-based access control** (Admin, Organization, Collaborator)
- **Comprehensive validation** with detailed error responses
- **Pagination support** for list endpoints
- **Rate limiting** and idempotency for critical operations

### Core Endpoints

#### Authentication
- `POST /v1/auth/signup` - Register new user
- `POST /v1/auth/login` - Login user
- `GET /v1/auth/me` - Get current user info

#### Offers (Collaborators)
- `POST /v1/offers` - Create donation offer
- `GET /v1/offers` - List offers (own offers for collaborators)
- `GET /v1/offers/:id` - Get offer details
- `PUT /v1/offers/:id` - Update offer

#### Offers (Organizations)
- `GET /v1/offers/nearby` - Get nearby offers with priority scores

#### Claims & Redemptions (Organizations)
- `POST /v1/claims` - Claim an offer
- `GET /v1/claims` - List organization's claims
- `POST /v1/redemptions` - Complete redemption (spend credits, award tokens)

#### Credits & Tokens
- `GET /v1/credits` - Get current month credits
- `GET /v1/tokens` - Get current month tokens

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DB_URL` | PostgreSQL connection string | `postgres://postgres:12345678@localhost:5432/foodflow?sslmode=disable` |
| `REDIS_URL` | Redis connection string | `redis://localhost:6379` |
| `JWT_SECRET` | JWT signing secret (change in production) | `your-super-secret-jwt-key-change-this-in-production` |
| `JWT_EXPIRATION` | JWT token expiration time | `24h` |
| `JWT_ISSUER` | JWT token issuer | `foodflow` |
| `SERVINGS_PER_CREDIT` | Servings per credit | `10` |
| `TOKEN_MULTIPLIER` | Token multiplier for rewards | `5` |
| `CLAIM_COOLOFF_SECONDS` | Cool-off period for auto-assignment | `600` |
| `RATE_LIMIT_RPS` | Rate limit requests per second | `10` |
| `RATE_LIMIT_BURST` | Rate limit burst capacity | `20` |

### Environment Configuration

**Important**: Update your `.env` file with proper values:
```bash
# Copy template and edit
cd foodflow-backend
cp env.sample .env
# Update JWT_SECRET with a secure random string in production
# Update DB_URL with your PostgreSQL credentials
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
- Tokens are tracked monthly and can be redeemed for rewards

## Development

### Available Commands

```bash
# Root level commands
make help                 # Show all available commands
make docker-compose-up   # Start PostgreSQL & Redis services

# Backend commands (run from foodflow-backend/ directory)
cd foodflow-backend
go build -o bin/foodflow cmd/api/main.go  # Build the application
go run cmd/api/main.go                    # Run locally (auto-loads .env)
go test ./...                             # Run tests
```

### Quick Start Commands
```bash
# Complete backend setup and start
cd foodflow-backend
cp env.sample .env        # Copy environment template
go mod tidy              # Install dependencies
cd ..                    # Go back to root
make docker-compose-up   # Start PostgreSQL & Redis
cd foodflow-backend      # Go to backend
go run cmd/api/main.go   # Start API (auto-loads .env)
```

### Database Migrations

**Auto-migration**: Migrations run automatically when the application starts.

### Backend Testing

```bash
cd foodflow-backend
go test ./...            # Run all tests
go test -cover ./...     # Run tests with coverage
```

## Monitoring

### Health Check
```bash
curl http://localhost:8080/health
```

### Logs
```bash
# Backend logs
cd foodflow-backend
docker-compose logs -f   # View service logs

# Service logs (from root)
make docker-compose-up   # Start services
docker-compose logs -f   # View all service logs
```

## Security Features

- **JWT Authentication**: HMAC-SHA256 signed tokens (simplified deployment)
- **Password Hashing**: Argon2id for secure password storage
- **Environment Loading**: Automatic .env file loading with godotenv
- **Rate Limiting**: IP and user-based rate limiting
- **Idempotency**: SHA256-based safe retry mechanisms
- **Input Validation**: Comprehensive request validation
- **SQL Injection Protection**: Parameterized queries
- **CORS**: Configurable cross-origin resource sharing

## Performance

- **Connection Pooling**: Optimized database connection management
- **Redis Caching**: Fast access to frequently used data
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

---

## Recent Updates

### v1.0 - Production Ready ✅
- **Complete API Implementation**: All 50+ endpoints implemented
- **JWT Authentication**: Switched to HMAC-SHA256 for simplified deployment
- **Environment Loading**: Added automatic .env file loading
- **Comprehensive Documentation**: Complete API documentation for frontend integration
- **Production Security**: SHA256 idempotency, rate limiting, input validation
- **Database Integration**: Full schema with auto-migration
- **Monorepo Structure**: Organized backend into foodflow-backend/ directory

### Backend Features Complete ✅
- ✅ User management with role-based access
- ✅ Organization and collaborator profiles
- ✅ Donation offer lifecycle management
- ✅ Smart matching algorithm with priority scoring
- ✅ Credit system with monthly allocation
- ✅ Token rewards for collaborators
- ✅ Claims and redemption processing
- ✅ Remote organization support
- ✅ Admin dashboard endpoints
- ✅ Comprehensive audit logging
- ✅ Rate limiting and idempotency
- ✅ Complete input validation
- ✅ Error handling and logging

**Ready for frontend integration** 🚀

## Next Steps

- 🚧 **Frontend Development**: React/Next.js application in `foodflow-frontend/`
- 📱 **Mobile App**: Flutter or React Native mobile application
- 🔧 **DevOps**: CI/CD pipelines and deployment automation
- 📊 **Analytics**: Advanced reporting and analytics dashboard
- 🔔 **Notifications**: Real-time notifications via WebSocket
- 🌐 **API Gateway**: API versioning and gateway implementation