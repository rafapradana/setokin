# Setokin Setup Script (Windows)

Write-Host "🚀 Setting up Setokin..." -ForegroundColor Cyan

# Copy environment file
if (-not (Test-Path ".env")) {
    Copy-Item ".env.example" ".env"
    Write-Host "✅ Created .env from .env.example" -ForegroundColor Green
    Write-Host "⚠️  Please update .env with your configuration" -ForegroundColor Yellow
} else {
    Write-Host "ℹ️  .env already exists, skipping" -ForegroundColor Blue
}

# Build and start services
Write-Host "🐳 Building Docker images..." -ForegroundColor Cyan
docker compose build

Write-Host "✅ Setup complete!" -ForegroundColor Green
Write-Host ""
Write-Host "Run 'make dev' to start the development environment"
