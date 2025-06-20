# Golang Fiber Backend with MongoDB and Refresh Tokens

A robust backend service built with Go Fiber framework, featuring MongoDB integration and JWT-based authentication with refresh token support.

## 🚀 Features

- **RESTful API** with Go Fiber framework
- **MongoDB** integration with proper indexing
- **JWT Authentication** with access and refresh tokens
- **Repository Pattern** for clean architecture
- **Input Validation** with comprehensive error handling
- **CORS** and security middleware
- **Docker** support for easy deployment
- **Refresh Token Rotation** for enhanced security
- **Password Reset via Email** with SendGrid integration
- **Swagger Documentation** with runtime enable/disable control
- **Environment-based Configuration** for development and production

## 📁 Project Structure

```
backend/
├── config/                    # Configuration management
├── models/                    # Data models and DTOs
├── repositories/              # Data access layer
│   ├── interfaces/           # Repository contracts
│   ├── user_repository.go    # User data operations
│   └── token_repository.go   # Token data operations
├── services/                  # Business logic layer
├── handlers/                  # HTTP request handlers
├── middleware/                # Custom middleware
├── routes/                    # Route definitions
├── database/                  # Database connection
├── utils/                     # Utility functions
├── docs/                      # Generated Swagger documentation
├── docker-compose.yml         # Docker development setup
├── Dockerfile                 # Container configuration
├── Makefile                   # Build and development commands
└── main.go                   # Application entry point
```

## 🛠️ Setup & Installation

### Prerequisites

- Go 1.21 or higher
- MongoDB (or Docker for containerized setup)

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd backend
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Setup environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Start MongoDB** (if not using Docker)
   ```bash
   # Start MongoDB locally or use the Docker option below
   ```

5. **Generate Swagger documentation and run**
   ```bash
   make dev
   # or manually:
   swag init -g main.go -o ./docs
   go run main.go
   ```

### Docker Development

1. **Start services with Docker Compose**
   ```bash
   docker-compose up -d
   ```

   This will start:
   - MongoDB on port 27017
   - Backend API on port 3000

2. **View logs**
   ```bash
   docker-compose logs -f backend
   ```

3. **Stop services**
   ```bash
   docker-compose down
   ```

## 📚 API Documentation

### Base URL
```
http://localhost:3000/api/v1
```

### API Documentation
- **Swagger UI**: `http://localhost:3000/swagger/index.html` *(when enabled)*
- **API Spec JSON**: `http://localhost:3000/swagger/doc.json`
- **Swagger Status**: `http://localhost:3000/api/v1/swagger-status`

### Health Check
```
GET /health
```

### Authentication Endpoints

#### Register User
```http
POST /auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

#### Login User
```http
POST /auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

#### Refresh Token
```http
POST /auth/refresh
Content-Type: application/json

{
  "refresh_token": "your-refresh-token"
}
```

#### Logout
```http
POST /auth/logout
Content-Type: application/json

{
  "refresh_token": "your-refresh-token"
}
```

#### Forgot Password
```http
POST /auth/forgot-password
Content-Type: application/json

{
  "email": "john@example.com"
}
```
*Generates a new secure password and sends it via email using SendGrid*

### Protected User Endpoints
*Requires Authorization header: `Bearer <access_token>`*

#### Get User Profile
```http
GET /users/profile
Authorization: Bearer <access_token>
```

#### Update User Profile
```http
PUT /users/profile
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "name": "Jane Doe",
  "email": "jane@example.com"
}
```

#### Delete User Profile
```http
DELETE /users/profile
Authorization: Bearer <access_token>
```

#### Logout from All Devices
```http
POST /users/logout-all
Authorization: Bearer <access_token>
```

## 🔐 Authentication Flow

1. **Register/Login** → Receive access token (15 min) + refresh token (7 days)
2. **API Requests** → Use access token in Authorization header
3. **Token Refresh** → Use refresh token to get new access token
4. **Token Rotation** → New refresh token provided on each refresh
5. **Logout** → Revoke specific refresh token
6. **Logout All** → Revoke all user's refresh tokens
7. **Password Reset** → Generate new password and send via email

## 🔧 Configuration

Environment variables in `.env`:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `3000` |
| `MONGODB_URI` | MongoDB connection string | `mongodb://localhost:27017/esp_backend_db` |
| `JWT_ACCESS_SECRET` | JWT access token secret | - |
| `JWT_REFRESH_SECRET` | JWT refresh token secret | - |
| `JWT_ACCESS_EXPIRY` | Access token expiry | `15m` |
| `JWT_REFRESH_EXPIRY` | Refresh token expiry | `168h` |
| `BCRYPT_ROUNDS` | Password hashing rounds | `12` |
| `SENDGRID_API_KEY` | SendGrid API key for email sending | - |
| `SENDGRID_FROM_EMAIL` | From email address for notifications | - |
| `SENDGRID_FROM_NAME` | From name for notifications | - |
| `RESET_PASSWORD_SUBJECT` | Subject line for password reset emails | `Password Reset - Your Account` |
| `PASSWORD_RESET_LENGTH` | Length of generated passwords | `10` |
| `PASSWORD_RESET_ATTEMPTS` | Max password reset attempts per hour | `3` |
| `SWAGGER_ENABLED` | Enable/disable Swagger UI | `true` (dev), `false` (prod) |
| `SWAGGER_HOST` | Swagger host for documentation | `localhost:3000` |
| `SWAGGER_BASE_PATH` | API base path | `/api/v1` |
| `SWAGGER_SCHEMES` | Supported schemes | `http` (dev), `https` (prod) |

## 📋 Swagger Documentation

### Runtime Configuration Control

Swagger can be enabled/disabled at runtime based on environment:

```bash
# Development (Swagger enabled by default)
APP_ENV=development SWAGGER_ENABLED=true go run main.go

# Production (Swagger disabled by default)
APP_ENV=production SWAGGER_ENABLED=false go run main.go

# Force enable in production (not recommended)
APP_ENV=production SWAGGER_ENABLED=true go run main.go
```

### Swagger Commands

```bash
# Generate Swagger documentation
make swagger-generate
# or: swag init -g main.go -o ./docs

# Generate docs and start server
make swagger-serve

# Validate Swagger spec
make swagger-validate

# Development mode (auto-generate docs)
make dev
```

### Environment-Specific Behavior

- **Development**: Swagger enabled by default, accessible at `/swagger/`
- **Staging**: Swagger enabled if explicitly configured
- **Production**: Swagger disabled by default for security

## 🧪 Testing

### Using Swagger UI (Recommended)

1. **Start the server**: `make dev`
2. **Open Swagger UI**: `http://localhost:3000/swagger/index.html`
3. **Test endpoints** directly from the browser interface
4. **View API documentation** with examples and schemas

### Manual Testing with curl

```bash
# Check Swagger status
curl -X GET http://localhost:3000/api/v1/swagger-status

# Register
curl -X POST http://localhost:3000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Forgot Password
curl -X POST http://localhost:3000/api/v1/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com"}'

# Get Profile (replace TOKEN with actual access token)
curl -X GET http://localhost:3000/api/v1/users/profile \
  -H "Authorization: Bearer TOKEN"
```

## 🏗️ Architecture

This project follows the **Repository Pattern** and **Clean Architecture** principles:

- **Models**: Define data structures and DTOs
- **Repository Layer**: Abstract data access with interfaces
- **Service Layer**: Contains business logic and use cases  
- **Handler Layer**: HTTP request/response handling
- **Middleware**: Cross-cutting concerns (auth, logging, CORS)

## 🔒 Security Features

- **Password Hashing** with bcrypt
- **JWT Access Tokens** (short-lived, 15 minutes)
- **Refresh Tokens** (longer-lived, 7 days)
- **Token Rotation** on refresh
- **Token Revocation** support
- **CORS** protection
- **Input Validation** with custom rules

## 📝 Response Format

All API responses follow this format:

```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... },
  "error": null
}
```

Error responses:
```json
{
  "success": false,
  "message": "Error description",
  "data": null,
  "error": "Detailed error message"
}
```

## 🚦 Status Codes

- `200` - Success
- `201` - Created
- `400` - Bad Request (validation errors)
- `401` - Unauthorized (invalid/expired token)
- `404` - Not Found
- `409` - Conflict (duplicate email)
- `500` - Internal Server Error

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License. 