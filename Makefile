.PHONY: dcup dcdown gup gdown gorun

DEFAULT: start

# docker-compose
dcup:
	docker-compose up -d

dcdown:
	docker-compose down

# goose migration
gup:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=postgres://user:qwerty@localhost:5432/postgres goose -dir migrations up
gdown:
	goose -dir migrations down

gorun:
	go run cmd/main.go

start: dcup wait-for-db gup gorun

wait-for-db:
	@echo "Waiting for database to be ready..."
	@until docker exec -it wb_task_one-postgres-1 pg_isready -U user; do \
		sleep 1; \
	done

pub:
	@json_data=$$(cat order.json); \
	nats pub example-subject "$$json_data"