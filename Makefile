# =========================
# Load env
# =========================
include .env
export

# =========================
# Config
# =========================
GOOSE=goose
MIGRATIONS_DIR=migrations
DB_DRIVER=postgres

# =========================
# Proto
# =========================
proto:
	buf generate

# =========================
# Database Migration
# =========================
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "❌ name is required"; \
		echo "👉 contoh:"; \
		echo "   make migrate-create name=create_users"; \
		exit 1; \
	fi
	$(GOOSE) -dir $(MIGRATIONS_DIR) create $(name) sql

migrate:
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_DSN)" up

migrate-down:
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_DSN)" down

migrate-status:
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_DSN)" status

db-reset:
	@echo "⚠️  RESET DATABASE (ALL MIGRATIONS WILL BE ROLLED BACK)"
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_DSN)" reset

db-refresh:
	@echo "🔥 RESET + MIGRATE DATABASE"
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_DSN)" reset
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_DSN)" up

# =========================
# Run Servers
# =========================
grpc:
	go run cmd/grpc/main.go

http:
	go run cmd/http/main.go

mail:
	go run cmd/mail_worker/main.go

run:
	make -j3 grpc http mail

module:
	@if [ -z "$(name)" ]; then \
		echo "❌ usage: make module name=module_name"; \
		exit 1; \
	fi
	@bash scripts/gen-module.sh $(name)


# =========================
# Helpers
# =========================
.PHONY: proto migrate migrate-down migrate-status migrate-create grpc http mail run db-reset db-refresh module	
