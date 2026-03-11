# Setokin Project Structure

**Version:** 1.0.0  
**Last Updated:** 11 Maret 2026

## Table of Contents

1. [Overview](#overview)
2. [Directory Structure](#directory-structure)
3. [Technology Stack](#technology-stack)
4. [Environment Setup](#environment-setup)
5. [Development Workflow](#development-workflow)
6. [Deployment Strategy](#deployment-strategy)
7. [Best Practices](#best-practices)

---

## Overview

Setokin menggunakan **monorepo architecture** dengan struktur yang terorganisir untuk memudahkan development, testing, dan deployment. Project ini dirancang agar developer baru dapat setup dengan cepat dan mulai berkontribusi dalam waktu singkat.

### Key Principles

- **Single source of truth** - Satu `.env` file di root
- **Docker-first development** - Semua service run via Docker Compose
- **Makefile automation** - Command sederhana untuk semua operasi
- **Clear separation** - Frontend, backend, dan infrastructure terpisah
- **Production-ready** - Support local dev, staging, dan production deployment

---

## Directory Structure

```
setokin/
├── .github/
│   └── workflows/
│       ├── ci.yml                    # CI pipeline (test, lint)
│       ├── deploy-staging.yml        # Deploy to staging
│       └── deploy-production.yml     # Deploy to production
│
├── api/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go              # Application entry point
│   ├── internal/
│   │   ├── config/
│   │   │   └── config.go            # Configuration management
│   │   ├── database/
│   │   │   ├── postgres.go          # PostgreSQL connection
│   │   │   └── migrations.go        # Migration runner
│   │   ├── middleware/
│   │   │   ├── auth.go              # JWT authentication
│   │   │   ├── cors.go              # CORS configuration
│   │   │   ├── logger.go            # Request logging
│   │   │   └── rate_limit.go        # Rate limiting
│   │   ├── models/
│   │   │   ├── user.go              # User model
│   │   │   ├── item.go              # Item model
│   │   │   ├── batch.go             # Batch model
│   │   │   ├── stock_in.go          # Stock In model
│   │   │   ├── stock_out.go         # Stock Out model
│   │   │   ├── category.go          # Category model
│   │   │   ├── unit.go              # Unit model
│   │   │   └── supplier.go          # Supplier model
│   │   ├── handlers/
│   │   │   ├── auth.go              # Auth endpoints
│   │   │   ├── users.go             # User endpoints
│   │   │   ├── items.go             # Item endpoints
│   │   │   ├── categories.go        # Category endpoints
│   │   │   ├── units.go             # Unit endpoints
│   │   │   ├── suppliers.go         # Supplier endpoints
│   │   │   ├── stock_in.go          # Stock In endpoints
│   │   │   ├── stock_out.go         # Stock Out endpoints
│   │   │   ├── batches.go           # Batch endpoints
│   │   │   ├── inventory.go         # Inventory endpoints
│   │   │   ├── reports.go           # Report endpoints
│   │   │   └── uploads.go           # File upload endpoints
│   │   ├── services/
│   │   │   ├── auth_service.go      # Auth business logic
│   │   │   ├── item_service.go      # Item business logic
│   │   │   ├── stock_service.go     # Stock management logic
│   │   │   ├── fefo_service.go      # FEFO algorithm implementation
│   │   │   ├── report_service.go    # Report generation
│   │   │   └── upload_service.go    # File upload with MinIO
│   │   ├── repositories/
│   │   │   ├── user_repo.go         # User data access
│   │   │   ├── item_repo.go         # Item data access
│   │   │   ├── batch_repo.go        # Batch data access
│   │   │   ├── stock_repo.go        # Stock transaction data access
│   │   │   └── report_repo.go       # Report data access
│   │   ├── utils/
│   │   │   ├── jwt.go               # JWT utilities
│   │   │   ├── hash.go              # Password hashing
│   │   │   ├── validator.go         # Input validation
│   │   │   ├── response.go          # Standard API response
│   │   │   └── pagination.go        # Pagination helpers
│   │   └── routes/
│   │       └── routes.go            # Route definitions
│   ├── pkg/
│   │   ├── logger/
│   │   │   └── logger.go            # Structured logging
│   │   └── errors/
│   │       └── errors.go            # Custom error types
│   ├── tests/
│   │   ├── integration/
│   │   │   ├── auth_test.go
│   │   │   ├── stock_test.go
│   │   │   └── fefo_test.go
│   │   └── unit/
│   │       ├── services/
│   │       └── repositories/
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   └── .air.toml                    # Hot reload config
│
├── web/
│   ├── public/
│   │   ├── favicon.ico
│   │   └── images/
│   ├── .env.local                   # Local env (gitignored)
│   ├── src/
│   │   ├── app/
│   │   │   ├── (auth)/
│   │   │   │   ├── login/
│   │   │   │   │   └── page.tsx
│   │   │   │   └── register/
│   │   │   │       └── page.tsx
│   │   │   ├── (dashboard)/
│   │   │   │   ├── layout.tsx       # Dashboard layout
│   │   │   │   ├── page.tsx         # Dashboard home
│   │   │   │   ├── inventory/
│   │   │   │   │   └── page.tsx
│   │   │   │   ├── items/
│   │   │   │   │   ├── page.tsx
│   │   │   │   │   ├── [id]/
│   │   │   │   │   │   └── page.tsx
│   │   │   │   │   └── new/
│   │   │   │   │       └── page.tsx
│   │   │   │   ├── stock-in/
│   │   │   │   │   ├── page.tsx
│   │   │   │   │   └── new/
│   │   │   │   │       └── page.tsx
│   │   │   │   ├── stock-out/
│   │   │   │   │   ├── page.tsx
│   │   │   │   │   └── new/
│   │   │   │   │       └── page.tsx
│   │   │   │   ├── batches/
│   │   │   │   │   └── page.tsx
│   │   │   │   ├── reports/
│   │   │   │   │   ├── page.tsx
│   │   │   │   │   ├── daily/
│   │   │   │   │   ├── weekly/
│   │   │   │   │   └── monthly/
│   │   │   │   ├── categories/
│   │   │   │   │   └── page.tsx
│   │   │   │   ├── suppliers/
│   │   │   │   │   └── page.tsx
│   │   │   │   └── settings/
│   │   │   │       └── page.tsx
│   │   │   ├── api/
│   │   │   │   └── [...proxy]/
│   │   │   │       └── route.ts     # API proxy to backend
│   │   │   ├── layout.tsx           # Root layout
│   │   │   └── page.tsx             # Landing page
│   │   ├── components/
│   │   │   ├── ui/                  # Shadcn UI components
│   │   │   │   ├── button.tsx
│   │   │   │   ├── input.tsx
│   │   │   │   ├── table.tsx
│   │   │   │   ├── dialog.tsx
│   │   │   │   ├── select.tsx
│   │   │   │   ├── card.tsx
│   │   │   │   └── ...
│   │   │   ├── layout/
│   │   │   │   ├── header.tsx
│   │   │   │   ├── sidebar.tsx
│   │   │   │   └── footer.tsx
│   │   │   ├── forms/
│   │   │   │   ├── item-form.tsx
│   │   │   │   ├── stock-in-form.tsx
│   │   │   │   ├── stock-out-form.tsx
│   │   │   │   └── ...
│   │   │   ├── tables/
│   │   │   │   ├── items-table.tsx
│   │   │   │   ├── batches-table.tsx
│   │   │   │   └── ...
│   │   │   └── charts/
│   │   │       ├── usage-chart.tsx
│   │   │       └── stock-chart.tsx
│   │   ├── lib/
│   │   │   ├── api.ts               # API client
│   │   │   ├── auth.ts              # Auth utilities
│   │   │   ├── utils.ts             # Helper functions
│   │   │   └── constants.ts         # App constants
│   │   ├── hooks/
│   │   │   ├── use-auth.ts          # Auth hook
│   │   │   ├── use-items.ts         # Items data hook
│   │   │   ├── use-stock.ts         # Stock data hook
│   │   │   └── use-reports.ts       # Reports data hook
│   │   ├── types/
│   │   │   ├── api.ts               # API types
│   │   │   ├── models.ts            # Data models
│   │   │   └── index.ts
│   │   └── styles/
│   │       └── globals.css
│   ├── next.config.js
│   ├── tailwind.config.ts
│   ├── tsconfig.json
│   ├── package.json
│   ├── Dockerfile
│   └── .eslintrc.json
│
├── mobile/
│   ├── android/
│   │   ├── app/
│   │   │   ├── src/
│   │   │   └── build.gradle
│   │   ├── gradle/
│   │   └── build.gradle
│   ├── ios/
│   │   ├── Runner/
│   │   ├── Runner.xcodeproj/
│   │   └── Runner.xcworkspace/
│   ├── lib/
│   │   ├── main.dart                # Application entry point
│   │   ├── app/
│   │   │   ├── routes/
│   │   │   │   └── app_routes.dart  # Route definitions
│   │   │   └── theme/
│   │   │       ├── app_theme.dart   # Theme configuration
│   │   │       └── app_colors.dart  # Color palette
│   │   ├── core/
│   │   │   ├── constants/
│   │   │   │   ├── api_constants.dart
│   │   │   │   └── app_constants.dart
│   │   │   ├── services/
│   │   │   │   ├── api_service.dart # HTTP client
│   │   │   │   ├── auth_service.dart
│   │   │   │   └── storage_service.dart
│   │   │   ├── utils/
│   │   │   │   ├── validators.dart
│   │   │   │   ├── formatters.dart
│   │   │   │   └── helpers.dart
│   │   │   └── errors/
│   │   │       └── app_exception.dart
│   │   ├── data/
│   │   │   ├── models/
│   │   │   │   ├── user.dart
│   │   │   │   ├── item.dart
│   │   │   │   ├── batch.dart
│   │   │   │   ├── stock_in.dart
│   │   │   │   ├── stock_out.dart
│   │   │   │   ├── category.dart
│   │   │   │   ├── unit.dart
│   │   │   │   └── supplier.dart
│   │   │   ├── repositories/
│   │   │   │   ├── auth_repository.dart
│   │   │   │   ├── item_repository.dart
│   │   │   │   ├── stock_repository.dart
│   │   │   │   └── report_repository.dart
│   │   │   └── providers/
│   │   │       ├── auth_provider.dart
│   │   │       ├── item_provider.dart
│   │   │       └── stock_provider.dart
│   │   ├── presentation/
│   │   │   ├── screens/
│   │   │   │   ├── auth/
│   │   │   │   │   ├── login_screen.dart
│   │   │   │   │   └── register_screen.dart
│   │   │   │   ├── dashboard/
│   │   │   │   │   └── dashboard_screen.dart
│   │   │   │   ├── inventory/
│   │   │   │   │   └── inventory_screen.dart
│   │   │   │   ├── items/
│   │   │   │   │   ├── items_list_screen.dart
│   │   │   │   │   ├── item_detail_screen.dart
│   │   │   │   │   └── item_form_screen.dart
│   │   │   │   ├── stock_in/
│   │   │   │   │   ├── stock_in_list_screen.dart
│   │   │   │   │   └── stock_in_form_screen.dart
│   │   │   │   ├── stock_out/
│   │   │   │   │   ├── stock_out_list_screen.dart
│   │   │   │   │   └── stock_out_form_screen.dart
│   │   │   │   ├── batches/
│   │   │   │   │   └── batches_screen.dart
│   │   │   │   ├── reports/
│   │   │   │   │   ├── reports_screen.dart
│   │   │   │   │   ├── daily_report_screen.dart
│   │   │   │   │   ├── weekly_report_screen.dart
│   │   │   │   │   └── monthly_report_screen.dart
│   │   │   │   └── settings/
│   │   │   │       └── settings_screen.dart
│   │   │   ├── widgets/
│   │   │   │   ├── common/
│   │   │   │   │   ├── app_button.dart
│   │   │   │   │   ├── app_text_field.dart
│   │   │   │   │   ├── app_card.dart
│   │   │   │   │   ├── loading_indicator.dart
│   │   │   │   │   └── error_widget.dart
│   │   │   │   ├── items/
│   │   │   │   │   ├── item_card.dart
│   │   │   │   │   └── item_list_tile.dart
│   │   │   │   ├── batches/
│   │   │   │   │   └── batch_card.dart
│   │   │   │   └── charts/
│   │   │   │       ├── usage_chart.dart
│   │   │   │       └── stock_chart.dart
│   │   │   └── controllers/
│   │   │       ├── auth_controller.dart
│   │   │       ├── item_controller.dart
│   │   │       ├── stock_controller.dart
│   │   │       └── report_controller.dart
│   │   └── generated/
│   │       └── l10n/                # Localization files
│   ├── test/
│   │   ├── unit/
│   │   ├── widget/
│   │   └── integration/
│   ├── assets/
│   │   ├── images/
│   │   ├── icons/
│   │   └── fonts/
│   ├── pubspec.yaml
│   ├── pubspec.lock
│   ├── analysis_options.yaml
│   └── README.md
│
├── db/
│   ├── db.sql                       # Main schema
│   ├── migrations/
│   │   ├── 001_initial_schema.sql
│   │   ├── 002_add_indexes.sql
│   │   └── ...
│   └── seeds/
│       ├── dev_seed.sql             # Development data
│       └── prod_seed.sql            # Production initial data
│
├── nginx/
│   ├── nginx.conf                   # Main nginx config
│   ├── conf.d/
│   │   ├── default.conf             # Default site config
│   │   ├── api.conf                 # API proxy config
│   │   └── ssl.conf                 # SSL configuration
│   └── Dockerfile
│
├── scripts/
│   ├── setup.sh                     # Initial setup script
│   ├── setup.ps1                    # Windows setup script
│   ├── dev.sh                       # Start dev environment
│   ├── dev.ps1                      # Windows dev script
│   ├── test.sh                      # Run tests
│   ├── deploy.sh                    # Deployment script
│   └── backup.sh                    # Database backup
│
├── docs/
│   ├── prd.md                       # Product requirements
│   ├── api/
│   │   └── api.md                   # API documentation
│   ├── project-structure.md         # This file
│   ├── development-guide.md         # Development guide
│   ├── deployment-guide.md          # Deployment guide
│   └── architecture.md              # Architecture decisions
│
├── .env                             # Environment variables (gitignored)
├── .env.example                     # Environment template
├── .gitignore
├── docker-compose.yml               # Development compose
├── docker-compose.prod.yml          # Production compose
├── Makefile                         # Task automation
├── README.md                        # Project overview
└── LICENSE

```

---


## Technology Stack

### API (Backend)

| Technology | Version | Purpose |
|------------|---------|---------|
| Go | 1.21+ | Backend language |
| Fiber | v2.52+ | Web framework |
| GORM | v1.25+ | ORM |
| JWT-Go | v5.0+ | JWT authentication |
| Air | latest | Hot reload |
| Testify | latest | Testing framework |

### Web (Frontend)

| Technology | Version | Purpose |
|------------|---------|---------|
| Next.js | 14+ | React framework |
| React | 18+ | UI library |
| TypeScript | 5+ | Type safety |
| Tailwind CSS | 3+ | Styling |
| Shadcn UI | latest | Component library |
| Phosphor Icons | latest | Icon set |
| React Query | latest | Data fetching |
| Zustand | latest | State management |

### Mobile

| Technology | Version | Purpose |
|------------|---------|---------|
| Flutter | 3.16+ | Mobile framework |
| Dart | 3.2+ | Programming language |
| GetX | 4.6+ | State management & routing |
| Dio | 5.4+ | HTTP client |
| Flutter Secure Storage | 9.0+ | Secure token storage |
| Freezed | 2.4+ | Code generation for models |
| Flutter Riverpod | 2.4+ | Alternative state management |
| FL Chart | 0.66+ | Charts and graphs |
| Cached Network Image | 3.3+ | Image caching |
| Image Picker | 1.0+ | Camera & gallery access |

### Infrastructure

| Technology | Version | Purpose |
|------------|---------|---------|
| PostgreSQL | 16+ | Database |
| MinIO | latest | Object storage |
| Nginx | 1.25+ | Reverse proxy |
| Docker | 24+ | Containerization |
| Docker Compose | 2.20+ | Orchestration |

---

## Environment Setup

### Prerequisites

- Docker & Docker Compose
- Make (GNU Make)
- Git

### Environment Variables

Single `.env` file di root directory:

```bash
# Application
APP_NAME=setokin
APP_ENV=development
APP_PORT=3000
API_PORT=8080

# Database
DB_HOST=postgres
DB_PORT=5432
DB_NAME=setokin
DB_USER=setokin
DB_PASSWORD=setokin_password
DB_SSL_MODE=disable

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# MinIO
MINIO_ENDPOINT=minio:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=setokin
MINIO_USE_SSL=false

# Web Frontend
NEXT_PUBLIC_API_URL=http://localhost:8080/v1
NEXT_PUBLIC_APP_URL=http://localhost:3000

# Mobile
FLUTTER_API_URL=http://10.0.2.2:8080/v1  # Android emulator
# FLUTTER_API_URL=http://localhost:8080/v1  # iOS simulator

# Nginx
NGINX_PORT=80
NGINX_SSL_PORT=443
```

### Quick Start

```bash
# 1. Clone repository
git clone https://github.com/your-org/setokin.git
cd setokin

# 2. Copy environment file
cp .env.example .env

# 3. Start development environment
make dev

# 4. Access application
# Web: http://localhost:3000
# API: http://localhost:8080
# MinIO Console: http://localhost:9001
# Mobile: Run from Android Studio / Xcode
```

---

## Development Workflow

### Makefile Commands

```makefile
# Development
make dev              # Start all services in development mode
make dev-web          # Start only web frontend
make dev-api          # Start only API backend
make dev-mobile       # Start mobile app (Flutter)
make dev-db           # Start only database

# Testing
make test             # Run all tests
make test-api         # Run API tests
make test-web         # Run web tests
make test-mobile      # Run mobile tests
make test-integration # Run integration tests

# Database
make db-migrate       # Run database migrations
make db-seed          # Seed database with test data
make db-reset         # Reset database (drop + migrate + seed)
make db-backup        # Backup database
make db-restore       # Restore database from backup

# Code Quality
make lint             # Run linters
make lint-api         # Lint API code
make lint-web         # Lint web code
make lint-mobile      # Lint mobile code
make format           # Format code
make format-api       # Format API code
make format-web       # Format web code
make format-mobile    # Format mobile code

# Build
make build            # Build all services
make build-api        # Build API
make build-web        # Build web
make build-mobile     # Build mobile (APK/IPA)
make build-nginx      # Build nginx

# Production Simulation
make prod-local       # Run production build locally
make prod-test        # Test production build

# Cleanup
make clean            # Stop and remove all containers
make clean-volumes    # Remove all volumes (WARNING: data loss)
make clean-all        # Clean everything

# Utilities
make logs             # View all logs
make logs-api         # View API logs
make logs-web         # View web logs
make logs-db          # View database logs
make shell-api        # Shell into API container
make shell-web        # Shell into web container
make shell-db         # Shell into database container
```

---


## Deployment Strategy

### Local Development

```yaml
# docker-compose.yml
version: '3.9'

services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: ${DB_NAME}
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./db/db.sql:/docker-entrypoint-initdb.d/01-schema.sql
      - ./db/seeds/dev_seed.sql:/docker-entrypoint-initdb.d/02-seed.sql
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5

  minio:
    image: minio/minio:latest
    command: server /data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: ${MINIO_ACCESS_KEY}
      MINIO_ROOT_PASSWORD: ${MINIO_SECRET_KEY}
    volumes:
      - minio_data:/data
    ports:
      - "9000:9000"
      - "9001:9001"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9000/minio/health/live"]
      interval: 30s
      timeout: 20s
      retries: 3

  backend:
    build:
      context: ./api
      dockerfile: Dockerfile
      target: development
    environment:
      - DB_HOST=${DB_HOST}
      - DB_PORT=${DB_PORT}
      - DB_NAME=${DB_NAME}
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - JWT_SECRET=${JWT_SECRET}
      - MINIO_ENDPOINT=${MINIO_ENDPOINT}
      - MINIO_ACCESS_KEY=${MINIO_ACCESS_KEY}
      - MINIO_SECRET_KEY=${MINIO_SECRET_KEY}
    volumes:
      - ./backend:/app
      - /app/tmp
    ports:
      - "${API_PORT}:8080"
    depends_on:
      postgres:
        condition: service_healthy
      minio:
        condition: service_healthy
    command: air -c .air.toml

  frontend:
    build:
      context: ./web
      dockerfile: Dockerfile
      target: development
    environment:
      - NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL}
      - NEXT_PUBLIC_APP_URL=${NEXT_PUBLIC_APP_URL}
    volumes:
      - ./frontend:/app
      - /app/node_modules
      - /app/.next
    ports:
      - "${APP_PORT}:3000"
    depends_on:
      - backend
    command: npm run dev

volumes:
  postgres_data:
  minio_data:
```

### Production Deployment (Ubuntu 24 LTS VPS)

#### Architecture

```
Internet
    ↓
[Nginx] (Port 80/443)
    ↓
    ├─→ [Web Container] (Next.js)
    ├─→ [API Container] (Go Fiber)
    ├─→ [PostgreSQL Container]
    └─→ [MinIO Container]

Mobile App (Flutter)
    ↓
[API Container] (Go Fiber)
```

#### GitHub Actions Workflow

```yaml
# .github/workflows/deploy-production.yml
name: Deploy to Production

on:
  push:
    branches: [main]
  workflow_dispatch:

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup SSH
        uses: webfactory/ssh-agent@v0.8.0
        with:
          ssh-private-key: ${{ secrets.SSH_PRIVATE_KEY }}

      - name: Deploy to VPS
        env:
          VPS_HOST: ${{ secrets.VPS_HOST }}
          VPS_USER: ${{ secrets.VPS_USER }}
        run: |
          ssh -o StrictHostKeyChecking=no $VPS_USER@$VPS_HOST << 'EOF'
            cd /opt/setokin
            git pull origin main
            make prod-deploy
          EOF

      - name: Health Check
        run: |
          sleep 30
          curl -f https://api.setokin.com/health || exit 1
```

#### Production Docker Compose

```yaml
# docker-compose.prod.yml
version: '3.9'

services:
  postgres:
    image: postgres:16-alpine
    restart: always
    environment:
      POSTGRES_DB: ${DB_NAME}
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./db/db.sql:/docker-entrypoint-initdb.d/01-schema.sql
      - ./backups:/backups
    networks:
      - setokin_network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5

  minio:
    image: minio/minio:latest
    restart: always
    command: server /data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: ${MINIO_ACCESS_KEY}
      MINIO_ROOT_PASSWORD: ${MINIO_SECRET_KEY}
    volumes:
      - minio_data:/data
    networks:
      - setokin_network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9000/minio/health/live"]
      interval: 30s
      timeout: 20s
      retries: 3

  backend:
    build:
      context: ./api
      dockerfile: Dockerfile
      target: production
    restart: always
    environment:
      - APP_ENV=production
      - DB_HOST=${DB_HOST}
      - DB_PORT=${DB_PORT}
      - DB_NAME=${DB_NAME}
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - JWT_SECRET=${JWT_SECRET}
      - MINIO_ENDPOINT=${MINIO_ENDPOINT}
      - MINIO_ACCESS_KEY=${MINIO_ACCESS_KEY}
      - MINIO_SECRET_KEY=${MINIO_SECRET_KEY}
    networks:
      - setokin_network
    depends_on:
      postgres:
        condition: service_healthy
      minio:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  frontend:
    build:
      context: ./web
      dockerfile: Dockerfile
      target: production
    restart: always
    environment:
      - NODE_ENV=production
      - NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL}
      - NEXT_PUBLIC_APP_URL=${NEXT_PUBLIC_APP_URL}
    networks:
      - setokin_network
    depends_on:
      - backend
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:3000"]
      interval: 30s
      timeout: 10s
      retries: 3

  nginx:
    build:
      context: ./nginx
      dockerfile: Dockerfile
    restart: always
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/conf.d:/etc/nginx/conf.d:ro
      - ./nginx/ssl:/etc/nginx/ssl:ro
      - nginx_logs:/var/log/nginx
    networks:
      - setokin_network
    depends_on:
      - frontend
      - backend
    healthcheck:
      test: ["CMD", "nginx", "-t"]
      interval: 30s
      timeout: 10s
      retries: 3

networks:
  setokin_network:
    driver: bridge

volumes:
  postgres_data:
  minio_data:
  nginx_logs:
```

---


## Best Practices

### API (Backend - Go)

#### Project Structure

```
internal/
├── config/      # Configuration management
├── database/    # Database connection & migrations
├── middleware/  # HTTP middleware
├── models/      # GORM models
├── handlers/    # HTTP handlers (thin layer)
├── services/    # Business logic (thick layer)
├── repositories/# Data access layer
├── utils/       # Utility functions
└── routes/      # Route definitions
```

#### Naming Conventions

- **Files**: `snake_case.go`
- **Packages**: `lowercase` (single word preferred)
- **Types**: `PascalCase`
- **Functions**: `PascalCase` (exported), `camelCase` (unexported)
- **Variables**: `camelCase`
- **Constants**: `PascalCase` or `SCREAMING_SNAKE_CASE`

#### Code Organization

```go
// handlers/item_handler.go
type ItemHandler struct {
    service services.ItemService
}

func NewItemHandler(service services.ItemService) *ItemHandler {
    return &ItemHandler{service: service}
}

func (h *ItemHandler) GetItems(c *fiber.Ctx) error {
    // Thin layer - just handle HTTP concerns
    items, err := h.service.GetItems(c.Context())
    if err != nil {
        return c.Status(500).JSON(utils.ErrorResponse(err))
    }
    return c.JSON(utils.SuccessResponse(items))
}

// services/item_service.go
type ItemService interface {
    GetItems(ctx context.Context) ([]models.Item, error)
}

type itemService struct {
    repo repositories.ItemRepository
}

func (s *itemService) GetItems(ctx context.Context) ([]models.Item, error) {
    // Thick layer - business logic here
    return s.repo.FindAll(ctx)
}

// repositories/item_repository.go
type ItemRepository interface {
    FindAll(ctx context.Context) ([]models.Item, error)
}

type itemRepository struct {
    db *gorm.DB
}

func (r *itemRepository) FindAll(ctx context.Context) ([]models.Item, error) {
    // Data access only
    var items []models.Item
    err := r.db.WithContext(ctx).
        Preload("Category").
        Preload("Unit").
        Where("is_active = ?", true).
        Find(&items).Error
    return items, err
}
```

#### Error Handling

```go
// pkg/errors/errors.go
type AppError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Status  int    `json:"-"`
}

var (
    ErrNotFound = &AppError{
        Code:    "resource_not_found",
        Message: "Resource not found",
        Status:  404,
    }
    ErrUnauthorized = &AppError{
        Code:    "unauthorized",
        Message: "Unauthorized access",
        Status:  401,
    }
)
```

#### Testing

```go
// services/item_service_test.go
func TestItemService_GetItems(t *testing.T) {
    // Arrange
    mockRepo := new(mocks.ItemRepository)
    service := NewItemService(mockRepo)
    
    expectedItems := []models.Item{
        {ID: uuid.New(), Name: "Test Item"},
    }
    mockRepo.On("FindAll", mock.Anything).Return(expectedItems, nil)
    
    // Act
    items, err := service.GetItems(context.Background())
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expectedItems, items)
    mockRepo.AssertExpectations(t)
}
```

---

### Web (Frontend - Next.js)

#### Project Structure

```
src/
├── app/           # Next.js 14 App Router
├── components/    # React components
├── lib/           # Utilities & configurations
├── hooks/         # Custom React hooks
├── types/         # TypeScript types
└── styles/        # Global styles
```

#### Naming Conventions

- **Files**: `kebab-case.tsx` for components, `camelCase.ts` for utilities
- **Components**: `PascalCase`
- **Functions**: `camelCase`
- **Types/Interfaces**: `PascalCase`
- **Constants**: `SCREAMING_SNAKE_CASE`

#### Component Structure

```tsx
// components/forms/item-form.tsx
import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useItems } from '@/hooks/use-items';

interface ItemFormProps {
  initialData?: Item;
  onSuccess?: () => void;
}

export function ItemForm({ initialData, onSuccess }: ItemFormProps) {
  const [formData, setFormData] = useState(initialData || {});
  const { createItem, updateItem } = useItems();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      if (initialData) {
        await updateItem(initialData.id, formData);
      } else {
        await createItem(formData);
      }
      onSuccess?.();
    } catch (error) {
      console.error('Failed to save item:', error);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <Input
        label="Item Name"
        value={formData.name}
        onChange={(e) => setFormData({ ...formData, name: e.target.value })}
        required
      />
      <Button type="submit">Save</Button>
    </form>
  );
}
```

#### API Client

```typescript
// lib/api.ts
import axios from 'axios';

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor for auth token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor for token refresh
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      
      try {
        const refreshToken = localStorage.getItem('refresh_token');
        const { data } = await axios.post('/auth/refresh', { refreshToken });
        
        localStorage.setItem('access_token', data.access_token);
        originalRequest.headers.Authorization = `Bearer ${data.access_token}`;
        
        return api(originalRequest);
      } catch (refreshError) {
        // Redirect to login
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }
    
    return Promise.reject(error);
  }
);

export default api;
```

#### Custom Hooks

```typescript
// hooks/use-items.ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import api from '@/lib/api';
import type { Item, CreateItemDto } from '@/types/models';

export function useItems() {
  const queryClient = useQueryClient();

  const { data: items, isLoading } = useQuery({
    queryKey: ['items'],
    queryFn: async () => {
      const { data } = await api.get<{ data: Item[] }>('/items');
      return data.data;
    },
  });

  const createItem = useMutation({
    mutationFn: async (dto: CreateItemDto) => {
      const { data } = await api.post('/items', dto);
      return data.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['items'] });
    },
  });

  return {
    items,
    isLoading,
    createItem: createItem.mutateAsync,
  };
}
```

---

### Mobile (Flutter)

#### Project Structure

```
lib/
├── main.dart          # Entry point
├── app/               # App-level configuration
├── core/              # Core utilities & services
├── data/              # Data layer (models, repos, providers)
└── presentation/      # UI layer (screens, widgets, controllers)
```

#### Naming Conventions

- **Files**: `snake_case.dart`
- **Classes**: `PascalCase`
- **Functions**: `camelCase`
- **Variables**: `camelCase`
- **Constants**: `lowerCamelCase` or `SCREAMING_SNAKE_CASE`
- **Private members**: `_prefixWithUnderscore`

#### Code Organization

```dart
// data/models/item.dart
import 'package:freezed_annotation/freezed_annotation.dart';

part 'item.freezed.dart';
part 'item.g.dart';

@freezed
class Item with _$Item {
  const factory Item({
    required String id,
    required String name,
    required String categoryId,
    required String unitId,
    required double minimumStock,
    String? description,
    @Default(true) bool isActive,
    required DateTime createdAt,
    required DateTime updatedAt,
  }) = _Item;

  factory Item.fromJson(Map<String, dynamic> json) => _$ItemFromJson(json);
}

// data/repositories/item_repository.dart
abstract class ItemRepository {
  Future<List<Item>> getItems();
  Future<Item> getItemById(String id);
  Future<Item> createItem(CreateItemDto dto);
  Future<Item> updateItem(String id, UpdateItemDto dto);
  Future<void> deleteItem(String id);
}

class ItemRepositoryImpl implements ItemRepository {
  final ApiService _apiService;

  ItemRepositoryImpl(this._apiService);

  @override
  Future<List<Item>> getItems() async {
    final response = await _apiService.get('/items');
    return (response.data['data'] as List)
        .map((json) => Item.fromJson(json))
        .toList();
  }
}

// presentation/controllers/item_controller.dart
class ItemController extends GetxController {
  final ItemRepository _repository;
  
  ItemController(this._repository);

  final items = <Item>[].obs;
  final isLoading = false.obs;
  final error = Rxn<String>();

  @override
  void onInit() {
    super.onInit();
    fetchItems();
  }

  Future<void> fetchItems() async {
    try {
      isLoading.value = true;
      error.value = null;
      items.value = await _repository.getItems();
    } catch (e) {
      error.value = e.toString();
    } finally {
      isLoading.value = false;
    }
  }

  Future<void> createItem(CreateItemDto dto) async {
    try {
      isLoading.value = true;
      final item = await _repository.createItem(dto);
      items.add(item);
      Get.back();
      Get.snackbar('Success', 'Item created successfully');
    } catch (e) {
      Get.snackbar('Error', e.toString());
    } finally {
      isLoading.value = false;
    }
  }
}

// presentation/screens/items/items_list_screen.dart
class ItemsListScreen extends StatelessWidget {
  const ItemsListScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final controller = Get.find<ItemController>();

    return Scaffold(
      appBar: AppBar(
        title: const Text('Items'),
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () => Get.toNamed(Routes.ITEM_FORM),
          ),
        ],
      ),
      body: Obx(() {
        if (controller.isLoading.value) {
          return const Center(child: CircularProgressIndicator());
        }

        if (controller.error.value != null) {
          return Center(
            child: Text('Error: ${controller.error.value}'),
          );
        }

        return ListView.builder(
          itemCount: controller.items.length,
          itemBuilder: (context, index) {
            final item = controller.items[index];
            return ItemListTile(
              item: item,
              onTap: () => Get.toNamed(
                Routes.ITEM_DETAIL,
                arguments: item.id,
              ),
            );
          },
        );
      }),
    );
  }
}
```

#### API Service

```dart
// core/services/api_service.dart
import 'package:dio/dio.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class ApiService {
  late final Dio _dio;
  final FlutterSecureStorage _storage;

  ApiService(this._storage) {
    _dio = Dio(BaseOptions(
      baseURL: const String.fromEnvironment('API_URL'),
      connectTimeout: const Duration(seconds: 30),
      receiveTimeout: const Duration(seconds: 30),
    ));

    _dio.interceptors.add(InterceptorsWrapper(
      onRequest: (options, handler) async {
        final token = await _storage.read(key: 'access_token');
        if (token != null) {
          options.headers['Authorization'] = 'Bearer $token';
        }
        return handler.next(options);
      },
      onError: (error, handler) async {
        if (error.response?.statusCode == 401) {
          // Try to refresh token
          final refreshed = await _refreshToken();
          if (refreshed) {
            // Retry original request
            return handler.resolve(await _retry(error.requestOptions));
          }
        }
        return handler.next(error);
      },
    ));
  }

  Future<bool> _refreshToken() async {
    try {
      final refreshToken = await _storage.read(key: 'refresh_token');
      if (refreshToken == null) return false;

      final response = await _dio.post('/auth/refresh', data: {
        'refresh_token': refreshToken,
      });

      await _storage.write(
        key: 'access_token',
        value: response.data['data']['access_token'],
      );
      await _storage.write(
        key: 'refresh_token',
        value: response.data['data']['refresh_token'],
      );

      return true;
    } catch (e) {
      return false;
    }
  }

  Future<Response> _retry(RequestOptions requestOptions) async {
    final options = Options(
      method: requestOptions.method,
      headers: requestOptions.headers,
    );
    return _dio.request(
      requestOptions.path,
      data: requestOptions.data,
      queryParameters: requestOptions.queryParameters,
      options: options,
    );
  }

  Future<Response> get(String path, {Map<String, dynamic>? queryParameters}) {
    return _dio.get(path, queryParameters: queryParameters);
  }

  Future<Response> post(String path, {dynamic data}) {
    return _dio.post(path, data: data);
  }

  Future<Response> put(String path, {dynamic data}) {
    return _dio.put(path, data: data);
  }

  Future<Response> delete(String path) {
    return _dio.delete(path);
  }
}
```

#### State Management with GetX

```dart
// app/routes/app_routes.dart
class Routes {
  static const LOGIN = '/login';
  static const DASHBOARD = '/dashboard';
  static const ITEMS = '/items';
  static const ITEM_DETAIL = '/items/:id';
  static const ITEM_FORM = '/items/form';
  static const STOCK_IN = '/stock-in';
  static const STOCK_OUT = '/stock-out';
  static const REPORTS = '/reports';
}

class AppPages {
  static final pages = [
    GetPage(
      name: Routes.LOGIN,
      page: () => const LoginScreen(),
      binding: AuthBinding(),
    ),
    GetPage(
      name: Routes.DASHBOARD,
      page: () => const DashboardScreen(),
      binding: DashboardBinding(),
    ),
    GetPage(
      name: Routes.ITEMS,
      page: () => const ItemsListScreen(),
      binding: ItemBinding(),
    ),
    // ... more routes
  ];
}

// Dependency injection
class ItemBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut(() => ItemController(Get.find<ItemRepository>()));
  }
}
```

#### Testing

```dart
// test/unit/controllers/item_controller_test.dart
void main() {
  late ItemController controller;
  late MockItemRepository mockRepository;

  setUp(() {
    mockRepository = MockItemRepository();
    controller = ItemController(mockRepository);
  });

  group('ItemController', () {
    test('fetchItems should update items list on success', () async {
      // Arrange
      final items = [
        Item(
          id: '1',
          name: 'Test Item',
          categoryId: 'cat1',
          unitId: 'unit1',
          minimumStock: 5.0,
          isActive: true,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        ),
      ];
      when(() => mockRepository.getItems()).thenAnswer((_) async => items);

      // Act
      await controller.fetchItems();

      // Assert
      expect(controller.items.length, 1);
      expect(controller.items.first.name, 'Test Item');
      expect(controller.isLoading.value, false);
      expect(controller.error.value, null);
    });

    test('fetchItems should set error on failure', () async {
      // Arrange
      when(() => mockRepository.getItems()).thenThrow(Exception('Network error'));

      // Act
      await controller.fetchItems();

      // Assert
      expect(controller.items.length, 0);
      expect(controller.error.value, isNotNull);
      expect(controller.isLoading.value, false);
    });
  });
}
```

---

### Database

#### Migration Strategy

```sql
-- db/migrations/001_initial_schema.sql
-- Always include rollback instructions

-- UP
CREATE SCHEMA IF NOT EXISTS data;
CREATE TABLE data.users (...);

-- DOWN
DROP TABLE IF EXISTS data.users;
DROP SCHEMA IF EXISTS data;
```

#### Query Optimization

```sql
-- Use indexes for frequent queries
CREATE INDEX batches_item_expiry_idx 
ON data.batches(item_id, expiry_date, is_depleted);

-- Use EXPLAIN ANALYZE to check query performance
EXPLAIN ANALYZE
SELECT * FROM data.batches 
WHERE item_id = $1 
  AND is_depleted = false 
ORDER BY expiry_date ASC;
```

---

### Docker

#### Multi-stage Builds

```dockerfile
# api/Dockerfile
# Development stage
FROM golang:1.21-alpine AS development
WORKDIR /app
RUN go install github.com/cosmtrek/air@latest
COPY go.mod go.sum ./
RUN go mod download
COPY . .
CMD ["air", "-c", ".air.toml"]

# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Production stage
FROM alpine:latest AS production
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

```dockerfile
# web/Dockerfile
# Development stage
FROM node:20-alpine AS development
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
EXPOSE 3000
CMD ["npm", "run", "dev"]

# Build stage
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# Production stage
FROM node:20-alpine AS production
WORKDIR /app
ENV NODE_ENV=production
COPY --from=builder /app/package*.json ./
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
RUN npm ci --only=production
EXPOSE 3000
CMD ["npm", "start"]
```

---

### Git Workflow

#### Branch Strategy

```
main (production)
  ↑
develop (staging)
  ↑
feature/xxx (feature branches)
```

#### Commit Convention

```
feat: add stock out FEFO logic
feat(mobile): add offline mode support
fix: resolve batch depletion bug
fix(api): fix JWT token refresh
docs: update API documentation
refactor: improve item service structure
test: add integration tests for stock operations
test(mobile): add widget tests
chore: update dependencies
chore(mobile): upgrade Flutter to 3.16
```

---

### Security

#### Environment Variables

```bash
# Never commit .env file
# Always use .env.example as template
# Rotate secrets regularly in production
```

#### API Security

```go
// Rate limiting
app.Use(limiter.New(limiter.Config{
    Max:        60,
    Expiration: 1 * time.Minute,
}))

// CORS
app.Use(cors.New(cors.Config{
    AllowOrigins: os.Getenv("ALLOWED_ORIGINS"),
    AllowMethods: "GET,POST,PUT,DELETE",
    AllowHeaders: "Origin,Content-Type,Authorization",
}))

// Helmet-like security headers
app.Use(func(c *fiber.Ctx) error {
    c.Set("X-Content-Type-Options", "nosniff")
    c.Set("X-Frame-Options", "DENY")
    c.Set("X-XSS-Protection", "1; mode=block")
    return c.Next()
})
```

---

## Monitoring & Logging

### Structured Logging

```go
// pkg/logger/logger.go
import "go.uber.org/zap"

var log *zap.Logger

func Init() {
    var err error
    if os.Getenv("APP_ENV") == "production" {
        log, err = zap.NewProduction()
    } else {
        log, err = zap.NewDevelopment()
    }
    if err != nil {
        panic(err)
    }
}

func Info(msg string, fields ...zap.Field) {
    log.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
    log.Error(msg, fields...)
}
```

### Health Checks

```go
// handlers/health.go
func HealthCheck(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "status": "healthy",
        "timestamp": time.Now(),
        "services": fiber.Map{
            "database": checkDatabase(),
            "minio": checkMinio(),
        },
    })
}
```

---

## Performance Optimization

### API (Backend)

- Use connection pooling for database
- Implement caching with Redis (future)
- Use goroutines for concurrent operations
- Optimize database queries with proper indexes
- Use pagination for large datasets

### Web (Frontend)

- Use Next.js Image optimization
- Implement code splitting
- Use React.memo for expensive components
- Lazy load routes and components
- Optimize bundle size

### Mobile

- Implement offline-first architecture
- Use cached_network_image for images
- Lazy load lists with pagination
- Optimize build size (remove unused packages)
- Use const constructors where possible
- Implement proper error boundaries

---

## Troubleshooting

### Common Issues

**Database connection failed**
```bash
# Check if postgres is running
make logs-db

# Reset database
make db-reset
```

**Port already in use**
```bash
# Find process using port
lsof -i :8080  # macOS/Linux
netstat -ano | findstr :8080  # Windows

# Kill process or change port in .env
```

**Docker build fails**
```bash
# Clean docker cache
docker system prune -a

# Rebuild without cache
make build --no-cache
```

---

## Additional Resources

- [Go Best Practices](https://go.dev/doc/effective_go)
- [Next.js Documentation](https://nextjs.org/docs)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [PostgreSQL Performance](https://www.postgresql.org/docs/current/performance-tips.html)

---

## Mobile Development Setup

### Prerequisites

- Flutter SDK 3.16+
- Android Studio (for Android development)
- Xcode (for iOS development, macOS only)
- Android SDK & Emulator
- iOS Simulator (macOS only)

### Setup Instructions

```bash
# 1. Install Flutter dependencies
cd mobile
flutter pub get

# 2. Generate code (Freezed, JSON serialization)
flutter pub run build_runner build --delete-conflicting-outputs

# 3. Run on Android emulator
flutter run

# 4. Run on iOS simulator (macOS only)
flutter run -d ios

# 5. Build APK (Android)
flutter build apk --release

# 6. Build IPA (iOS, macOS only)
flutter build ios --release
```

### Environment Configuration

```dart
// lib/core/constants/api_constants.dart
class ApiConstants {
  static const String baseUrl = String.fromEnvironment(
    'API_URL',
    defaultValue: 'http://10.0.2.2:8080/v1', // Android emulator
  );
  
  static const Duration connectTimeout = Duration(seconds: 30);
  static const Duration receiveTimeout = Duration(seconds: 30);
}
```

### Running with different environments

```bash
# Development
flutter run --dart-define=API_URL=http://10.0.2.2:8080/v1

# Staging
flutter run --dart-define=API_URL=https://staging-api.setokin.com/v1

# Production
flutter run --dart-define=API_URL=https://api.setokin.com/v1
```

### Mobile-specific Features

- **Offline Mode**: Cache data locally using Hive/SQLite
- **Camera Integration**: Scan barcodes for items
- **Push Notifications**: Expiry alerts via FCM
- **Biometric Auth**: Fingerprint/Face ID login
- **Dark Mode**: Support system theme
- **Localization**: Indonesian & English support

---

## Platform-specific Notes

### Android

- Minimum SDK: 21 (Android 5.0)
- Target SDK: 34 (Android 14)
- Permissions required:
  - Internet
  - Camera (for barcode scanning)
  - Storage (for file uploads)

### iOS

- Minimum iOS: 12.0
- Permissions required (Info.plist):
  - NSCameraUsageDescription
  - NSPhotoLibraryUsageDescription
  - NSLocationWhenInUseUsageDescription (future)

---

## Cross-platform Considerations

### API Compatibility

All three platforms (Web, Mobile) consume the same REST API:

```
Web (Next.js) ──┐
                ├──→ API (Go Fiber) ──→ PostgreSQL
Mobile (Flutter)┘
```

### Shared Features

- JWT authentication with refresh tokens
- FEFO stock management
- Real-time inventory updates
- Expiry alerts
- Comprehensive reporting

### Platform-specific Features

| Feature | Web | Mobile |
|---------|-----|--------|
| Offline Mode | ❌ | ✅ |
| Camera/Barcode | ❌ | ✅ |
| Push Notifications | ❌ | ✅ |
| File Upload | ✅ | ✅ |
| Responsive Design | ✅ | Native |
| PWA Support | ✅ | N/A |

---

**Last Updated:** 11 Maret 2026
**Maintained By:** Setokin Developer (Rafa)
