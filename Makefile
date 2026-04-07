.PHONY: help \
	backend-build backend-run backend-test backend-lint backend-fmt backend-clean \
	frontend-install frontend-dev frontend-build frontend-lint frontend-fmt frontend-clean \
	docker-up docker-down docker-build docker-clean \
	docker-logs-backend docker-logs-frontend docker-logs-db \
	docker-restart-backend docker-restart-frontend

-include .env
-include backend/.env.test
export

# ===========================================
# Help
# ===========================================

help: ## Display this help message
	@echo "Available commands:"
	@grep -hE '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2}'

# ===========================================
# Backend
# ===========================================

backend-build: ## Build backend binary
	@echo "Building backend..."
	@cd backend && go build -o ../bin/server ./cmd/server
	@echo "Binary: bin/server"

backend-run: ## Run backend locally (applies migrations on start)
	@cd backend && set -a && [ -f .env ] && . ./.env; set +a; go run ./cmd/server

backend-test: ## Run backend tests
	@echo "Running backend tests..."
	@cd backend && TEST_DATABASE_URL="$(TEST_DATABASE_URL)" go test -v -race -p 1 ./...

backend-lint: ## Lint backend code
	@cd backend && golangci-lint run ./...

backend-fmt: ## Format backend code
	@cd backend && gofmt -s -w .

backend-swagger: ## Generate swagger docs (requires swag CLI)
	@cd backend && swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal
	@echo "Docs: http://localhost:8080/docs"

backend-clean: ## Remove backend build artifacts
	@rm -f bin/server
	@rm -rf backend/docs

# ===========================================
# Frontend
# ===========================================

frontend-install: ## Install frontend dependencies
	@cd frontend && pnpm install

frontend-dev: ## Start frontend dev server
	@cd frontend && pnpm dev

frontend-build: ## Build frontend for production
	@cd frontend && pnpm build

frontend-lint: ## Lint frontend code
	@cd frontend && pnpm lint

frontend-fmt: ## Format frontend code
	@cd frontend && pnpm exec prettier --write src

frontend-test: ## Run frontend tests
	@cd frontend && pnpm test

frontend-test-coverage: ## Run frontend tests with coverage
	@cd frontend && pnpm test:coverage

frontend-clean: ## Remove frontend build artifacts
	@rm -rf frontend/dist

# ===========================================
# Docker
# ===========================================

docker-up: ## Start all services (db + backend + frontend)
	@echo "Starting services..."
	docker compose up -d

docker-build: ## Build and start all services
	@echo "Building and starting services..."
	docker compose up -d --build

docker-down: ## Stop all services
	docker compose down

docker-clean: ## Stop all services and remove volumes
	docker compose down -v

docker-logs-backend: ## Tail backend logs
	docker compose logs -f backend

docker-logs-frontend: ## Tail frontend logs
	docker compose logs -f frontend

docker-logs-db: ## Tail database logs
	docker compose logs -f db

docker-restart-backend: ## Restart backend service
	docker compose restart backend

docker-restart-frontend: ## Restart frontend service
	docker compose restart frontend
