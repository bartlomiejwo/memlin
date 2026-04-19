.PHONY: help backend-up backend-down backend-build backend-run frontend-up frontend-down frontend-build frontend-run all-up all-down migrate-up migrate-down migrate-force seed-db sqlc
.DEFAULT_GOAL := run

# 🆘 Help: Show available commands
help:
	@echo "Available commands:"
	@echo ""
	@echo " Backend:"
	@echo " make backend-up - Start backend containers"
	@echo " make backend-down - Stop backend containers"
	@echo " make backend-build - Build backend image"
	@echo " make backend-run - Run backend in detached mode"
	@echo " make migrate-up - Apply database migrations"
	@echo " make migrate-down - Rollback last migration"
	@echo " make migrate-force - Force apply migrations"
	@echo " make seed-db - Seed database"
	@echo " make sqlc - Generate Go code from SQL queries using sqlc"
	@echo ""
	@echo " Frontend:"
	@echo " make frontend-up - Start frontend containers"
	@echo " make frontend-down - Stop frontend containers"
	@echo " make frontend-build - Build frontend image"
	@echo " make frontend-run - Run frontend in detached mode"
	@echo ""
	@echo " General:"
	@echo " make run - Start backend + frontend"
	@echo " make down - Stop backend + frontend"

# 🚀 Backend Tasks
backend-up:
	docker compose up -d backend db pgadmin

backend-down:
	docker compose down backend db pgadmin

backend-build:
	docker compose build backend

backend-run:
	docker compose up --build backend db pgadmin

# 📦 Database Migrations
migrate-up:
	docker compose run --rm backend /app/cmd/migrate/main up

migrate-down:
	docker compose run --rm backend /app/cmd/migrate/main down

# 🧬 Code Generation
sqlc:
	cd backend && sqlc generate

# 🎨 Frontend Tasks
frontend-up:
	docker compose up -d frontend

frontend-down:
	docker compose down frontend

frontend-build:
	docker compose build frontend

frontend-run:
	docker compose up --build frontend

# 🏗️ Full Project
run:
	docker compose up --build

down:
	docker compose down