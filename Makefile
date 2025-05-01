all: migrate seed
	docker compose up -d

test: migrate_test
	docker compose run --rm test

migrate:
	docker compose run --rm app bash -c "goose -dir goose -allow-missing postgres 'postgres://user:pass@db/db?sslmode=disable' up"

seed:
	docker compose run --rm app bash -c "goose -dir dev/seeds postgres 'postgres://user:pass@db/db?sslmode=disable' up"

migrate_test:
	docker compose up -d test_db
	docker compose run --rm app bash -c "goose -dir goose postgres 'postgres://user:pass@test_db/db?sslmode=disable' up"

new_migration:
	docker compose run --rm app bash -c "goose -dir goose postgres 'postgres://user:pass@db/db?sslmode=disable' create change_me sql"
	