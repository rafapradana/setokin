#!/bin/bash
# Setokin Setup Script

set -e

echo "🚀 Setting up Setokin..."

# Copy environment file
if [ ! -f .env ]; then
    cp .env.example .env
    echo "✅ Created .env from .env.example"
    echo "⚠️  Please update .env with your configuration"
else
    echo "ℹ️  .env already exists, skipping"
fi

# Build and start services
echo "🐳 Building Docker images..."
docker compose build

echo "✅ Setup complete!"
echo ""
echo "Run 'make dev' to start the development environment"
