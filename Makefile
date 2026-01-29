.PHONY: help setup up stop logs db-create-migration db-migrate db-rollback

# Default target
help: ## Show available commands
	@echo "🚀 Website Dummy App - Development Commands (Local Environment)"
	@echo ""
	@echo "📋 Main Commands:"
	@echo "   make setup     - First time setup (migrations, seed data, etc.)"
	@echo "   make up        - Start development session (packages, migrations, containers)"
	@echo "   make stop      - Stop all services (non-destructive)"
	@echo ""
	@echo "🗄️  Database Commands:"
	@echo "   make db-create-migration name=migration_name  - Create new migration"
	@echo "   make db-migrate                               - Run pending migrations"
	@echo "   make db-rollback                              - Rollback last migration"
	@echo ""
	@echo "🌐 Your app will be available at:"
	@echo "   Frontend: http://localhost:3000"
	@echo "   Backend:  http://localhost:8080"
	@echo "   MailHog:  http://localhost:8025 (for viewing emails in local dev)"

# First time setup - run once when cloning repo
setup: ## First time setup (migrations, seed data, environment)
	@echo "🔧 First time setup..."
	@if [ ! -f .env ]; then \
		echo "📄 Copying env.example to .env..."; \
		cp env.example .env; \
		echo "✅ Environment file created"; \
	else \
		echo "✅ Environment file already exists"; \
	fi
	@echo "🐳 Starting database..."
	@docker-compose -f docker-compose.local.yml up -d postgres
	@echo "⏳ Waiting for database to be ready..."
	@sleep 10
	@echo "🗄️  Creating database..."
	@docker-compose -f docker-compose.local.yml exec postgres psql -U postgres -c "CREATE DATABASE \"website-dummy\";" || echo "Database may already exist"
	@echo "🗄️  Running migrations..."
	@$(MAKE) db-migrate
	@echo "🌱 Seeding data..."
	@# TODO: Add seed data command when we have it
	@echo "✅ Setup complete! Run 'make up' to start development."

# Daily development startup
up: ## Start development session (packages, migrations, containers)
	@echo "🚀 Starting development session..."
	@echo "📦 Checking Go dependencies..."
	@if command -v go >/dev/null 2>&1; then \
		cd backend && go mod download; \
		echo "✅ Go dependencies ready"; \
	else \
		echo "⚠️  Go not installed locally (will use Docker)"; \
	fi
	@echo "🐳 Starting database..."
	@docker-compose -f docker-compose.local.yml up -d postgres
	@echo "⏳ Waiting for database to be ready..."
	@sleep 10
	@echo "🗄️  Ensuring database exists..."
	@docker-compose -f docker-compose.local.yml exec postgres psql -U postgres -c "CREATE DATABASE \"website-dummy\";" || echo "Database already exists"
	@echo "🗄️  Running any pending migrations..."
	@$(MAKE) db-migrate
	@echo "🐳 Starting all services..."
	@docker-compose -f docker-compose.local.yml up --build -d
	@echo "📦 Ensuring Go dependencies are tidy..."
	@sleep 5
	@docker-compose -f docker-compose.local.yml exec backend go mod tidy || echo "Backend not ready yet, dependencies will be resolved on next restart"
	@echo ""
	@echo "🎉 All services started successfully!"
	@echo ""
	@echo "🌐 Your application is now running:"
	@echo "   Frontend:  http://localhost:3000"
	@echo "   Backend:   http://localhost:8080"
	@echo "   API Health: http://localhost:8080/api/health"
	@echo "   MailHog:   http://localhost:8025 (for viewing emails in local dev)"
	@echo ""
	@echo "📋 Useful commands:"
	@echo "   make logs      - View all logs"
	@echo "   make stop      - Stop all services"
	@echo "   make help      - See all commands"

# Stop services (non-destructive)
stop: ## Stop all services (non-destructive)
	@echo "🛑 Stopping all services..."
	@docker-compose -f docker-compose.local.yml stop
	@echo "✅ All services stopped. Data preserved."

# View logs
logs: ## View logs from all services
	@docker-compose -f docker-compose.local.yml logs -f

logs-backend: ## View logs from backend service only
	@echo "📋 Viewing backend logs..."
	@docker-compose -f docker-compose.local.yml logs -f backend

# Database migration commands
db-create-migration: ## Create new migration (usage: make db-create-migration name=create_users)
	@if [ -z "$(name)" ]; then \
		echo "❌ Please provide migration name: make db-create-migration name=migration_name"; \
		exit 1; \
	fi
	@echo "📝 Creating migration: $(name)"
	@if command -v migrate >/dev/null 2>&1; then \
		migrate create -ext sql -dir backend/migrations $(name); \
	else \
		docker run --rm -v $(PWD)/backend/migrations:/migrations migrate/migrate create -ext sql -dir /migrations $(name); \
	fi
	@echo "✅ Migration created in backend/migrations/ directory"

db-migrate: ## Run pending migrations
	@echo "🗄️  Running database migrations..."
	@if [ ! -f .env ]; then \
		echo "❌ .env file not found. Run 'make setup' first."; \
		exit 1; \
	fi
	@if [ ! -d backend/migrations ] || [ -z "$$(ls -A backend/migrations/*.sql 2>/dev/null)" ]; then \
		echo "ℹ️  No migration files found, skipping migrations"; \
	else \
		export $$(grep -v '^#' .env | sed 's/#.*//' | sed 's/[[:space:]]*$$//' | grep -v '^$$' | xargs) && \
		if command -v migrate >/dev/null 2>&1; then \
			migrate -path backend/migrations -database "postgres://$${DB_USER:-postgres}:$${DB_PASSWORD:-postgres}@$${DB_HOST:-localhost}:$${DB_PORT:-5432}/$${DB_NAME:-website-dummy}?sslmode=$${DB_SSLMODE:-disable}" up; \
		else \
			docker run --rm -v $(PWD)/backend/migrations:/migrations --network host migrate/migrate \
			-path /migrations -database "postgres://postgres:postgres@localhost:5432/website-dummy?sslmode=disable" up; \
		fi; \
	fi
	@echo "✅ Migrations complete"

db-rollback: ## Rollback last migration
	@echo "⏪ Rolling back last migration..."
	@if [ ! -f .env ]; then \
		echo "❌ .env file not found. Run 'make setup' first."; \
		exit 1; \
	fi
	@export $$(grep -v '^#' .env | sed 's/#.*//' | sed 's/[[:space:]]*$$//' | grep -v '^$$' | xargs) && \
	if command -v migrate >/dev/null 2>&1; then \
		migrate -path backend/migrations -database "postgres://$${DB_USER:-postgres}:$${DB_PASSWORD:-postgres}@$${DB_HOST:-localhost}:$${DB_PORT:-5432}/$${DB_NAME:-website-dummy}?sslmode=$${DB_SSLMODE:-disable}" down 1; \
	else \
		docker run --rm -v $(PWD)/backend/migrations:/migrations --network host migrate/migrate \
		-path /migrations -database "postgres://postgres:postgres@localhost:5432/website-dummy?sslmode=disable" down 1; \
	fi
	@echo "✅ Rollback complete"