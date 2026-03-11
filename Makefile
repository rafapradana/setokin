# =============================================================================
# Setokin Makefile
# =============================================================================

.PHONY: help dev dev-api dev-web dev-db test test-api test-web build clean \
        db-migrate db-seed db-reset db-backup lint format logs

# Default
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# =============================================================================
# Development
# =============================================================================

dev: ## Start all services in development mode
	docker compose up --build

dev-api: ## Start only API backend
	docker compose up --build backend postgres minio

dev-web: ## Start only web frontend
	docker compose up --build frontend

dev-db: ## Start only database
	docker compose up --build postgres

# =============================================================================
# Testing
# =============================================================================

test: test-api ## Run all tests

test-api: ## Run API tests
	cd api && go test -v -race ./...

test-web: ## Run web tests
	cd web && npm test

test-integration: ## Run integration tests
	cd api && go test -v -race ./tests/integration/...

# =============================================================================
# Database
# =============================================================================

db-migrate: ## Run database migrations
	docker compose exec postgres psql -U $${DB_USER:-setokin} -d $${DB_NAME:-setokin} -f /docker-entrypoint-initdb.d/01-schema.sql

db-seed: ## Seed database with test data
	docker compose exec postgres psql -U $${DB_USER:-setokin} -d $${DB_NAME:-setokin} -f /docker-entrypoint-initdb.d/02-seed.sql

db-reset: ## Reset database (drop + recreate + migrate + seed)
	docker compose down -v
	docker compose up -d postgres
	@echo "Waiting for postgres to start..."
	@sleep 5
	docker compose up -d

db-backup: ## Backup database
	@mkdir -p backups
	docker compose exec postgres pg_dump -U $${DB_USER:-setokin} $${DB_NAME:-setokin} > backups/backup_$$(date +%Y%m%d_%H%M%S).sql
	@echo "Backup created in backups/"

db-restore: ## Restore database from backup (usage: make db-restore FILE=backups/xxx.sql)
	docker compose exec -T postgres psql -U $${DB_USER:-setokin} $${DB_NAME:-setokin} < $(FILE)

# =============================================================================
# Code Quality
# =============================================================================

lint: lint-api ## Run all linters

lint-api: ## Lint API code
	cd api && go vet ./...

format: format-api format-web ## Format all code

format-api: ## Format API code
	cd api && gofmt -w .

format-web: ## Format web code
	cd web && npm run format

# =============================================================================
# Build
# =============================================================================

build: ## Build all services for production
	docker compose -f docker-compose.prod.yml build

build-api: ## Build API only
	docker compose -f docker-compose.prod.yml build backend

build-web: ## Build web only
	docker compose -f docker-compose.prod.yml build frontend

build-nginx: ## Build nginx only
	docker compose -f docker-compose.prod.yml build nginx

# =============================================================================
# Production
# =============================================================================

prod-local: ## Run production build locally
	docker compose -f docker-compose.prod.yml up --build

prod-deploy: ## Deploy with production compose
	docker compose -f docker-compose.prod.yml up -d --build

# =============================================================================
# Logs
# =============================================================================

logs: ## View all logs
	docker compose logs -f

logs-api: ## View API logs
	docker compose logs -f backend

logs-web: ## View web logs
	docker compose logs -f frontend

logs-db: ## View database logs
	docker compose logs -f postgres

# =============================================================================
# Shell Access
# =============================================================================

shell-api: ## Shell into API container
	docker compose exec backend sh

shell-web: ## Shell into web container
	docker compose exec frontend sh

shell-db: ## Shell into database container
	docker compose exec postgres psql -U $${DB_USER:-setokin} $${DB_NAME:-setokin}

# =============================================================================
# Cleanup
# =============================================================================

clean: ## Stop and remove all containers
	docker compose down

clean-volumes: ## Remove all volumes (WARNING: data loss)
	docker compose down -v

clean-all: ## Clean everything including images
	docker compose down -v --rmi all
