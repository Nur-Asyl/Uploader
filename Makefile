.PHONY: migrate_init migrate_up_postgres migrate_down_postgres migrate_fix_postgres start

migrate_init:
	migrate create -ext sql -dir migrations/postgres/ -seq excel_mg

migrate_up_postgres:
	migrate -path migrations/postgres/ -database "postgresql://postgres:root@localhost:5432/parser?sslmode=disable" -verbose up

migrate_down_postgres:
	migrate -path migrations/postgres/ -database "postgresql://postgres:root@localhost:5432/parser?sslmode=disable" -verbose down

migrate_fix_postgres:
	migrate -path migrations/postgres/ -database "postgresql://postgres:root@localhost:5432/parser?sslmode=disable" force 000001

start:
	go run cmd/app/main.go
