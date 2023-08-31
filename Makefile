include .env
export

compose-up:
	docker-compose up --build -d && docker-compose logs -f
.PHONY: compose-up

compose-down:
	docker-compose down --remove-orphans
.PHONY: compose-down

docker-rm-volume:
	docker volume rm pg-data
.PHONY: docker-rm-volume

migrate-up:
	migrate -path migrations -database '$(PG_URL)' up
.PHONY: migrate-up

migrate-down:
	echo "y" | migrate -path migrations -database '$(PG_URL)' down
.PHONY: migrate-down

swagger:
	swag init -g internal/app/app.go --parseInternal --parseDependency
.PHONY: swagger