# Manage App

A full-stack application with Go + Gin + GORM backend and Vue.js + Tailwind frontend, following clean architecture principles.

## 🏗️ Project Structure

```
manage/
├── main.go                          # Application entry point
├── internal/                        # Private application code
│   ├── config/                      # Configuration (DB, CORS)
│   │   ├── database.go
│   │   └── cors.go
│   ├── controllers/                 # HTTP handlers
│   │   ├── user_controller.go
│   │   └── health_controller.go
│   ├── services/                    # Business logic
│   │   └── user_service.go
│   ├── models/                      # Data models
│   │   └── user.go
│   ├── routes/                      # Route definitions
│   │   └── routes.go
│   └── middleware/                  # Custom middleware
│       └── logging.go
└── frontend/                        # Vue.js frontend
    ├── src/
    │   ├── App.vue
    │   ├── main.js
    │   └── style.css
    └── package.json
```

## 🚀 Quick Start

### First Time Setup (New Developers)
```bash
# Clone the repo and run setup once:
make setup
```

This will:
- ✅ Copy `env.example` to `.env`
- ✅ Start database
- ✅ Run migrations
- ✅ Seed initial data (when available)

### Daily Development
```bash
# Start your development session:
make up

# When done for the day:
make stop
```

### Access Your Application
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080  
- **Health Check**: http://localhost:8080/api/health
- **API Documentation**: [API.md](./API.md)

## 🔧 Development

### Hot Reload
Both frontend and backend support hot reload:
- **Backend**: Uses Air for Go hot reload (watches `internal/` directory)
- **Frontend**: Uses Vite dev server with hot reload
- **Database**: PostgreSQL with persistent data

### Making Changes
- Edit Go files in `internal/` - backend will restart automatically
- Edit Vue files in `frontend/src/` - frontend will hot reload
- Database changes persist in Docker volume
- Use migrations for database schema changes

### Database Migrations
```bash
# Create a new migration
make db-create-migration name=create_users_table

# Run pending migrations
make db-migrate

# Rollback last migration
make db-rollback
```

### Services
- **postgres**: PostgreSQL database on port 5432
- **backend**: Go API server on port 8080
- **frontend**: Vue.js app on port 3000

## 📚 API Architecture

### Clean Architecture Layers

1. **Controllers** (`internal/controllers/`)
   - Handle HTTP requests/responses
   - Input validation and serialization
   - Delegate business logic to services

2. **Services** (`internal/services/`)
   - Business logic and validation
   - Orchestrate data operations
   - Independent of HTTP layer

3. **Models** (`internal/models/`)
   - Data structures and database schema
   - GORM model definitions

4. **Config** (`internal/config/`)
   - Database connections
   - CORS configuration
   - Application setup

### API Endpoints

- `GET /api/health` - Health check with database status
- `GET /api/name` - Returns "Tyler" (legacy endpoint)
- `GET /api/users` - Get all users
- `POST /api/users` - Create new user
- `GET /api/users/:id` - Get user by ID

See [API.md](./API.md) for detailed documentation.

## 🛠️ Tech Stack

- **Backend**: Go 1.21, Gin, GORM, PostgreSQL
- **Frontend**: Vue.js 3, Tailwind CSS v4.1, Vite
- **Infrastructure**: Docker, Docker Compose
- **Development**: Air (Go hot reload), Vite HMR

## 🎯 Architecture Benefits

- **Separation of Concerns**: Controllers, services, and models are clearly separated
- **Testability**: Business logic in services can be easily unit tested
- **Maintainability**: Clean structure makes it easy to add new features
- **Scalability**: Easy to add new controllers, services, and models
- **Hot Reload**: Fast development cycle with automatic restarts
