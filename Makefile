migrate-up:
	docker run -v "${CURDIR}/migrations:/migrations" --network host migrate/migrate -path=/migrations/ -database "mysql://root:test100500@tcp(localhost:3306)/flat-search" up

migrate-down:
	docker run -v "${CURDIR}/migrations:/migrations" --network host migrate/migrate -path=/migrations/ -database "mysql://root:test100500@tcp(localhost:3306)/flat-search" down -all

lint:
	docker-compose run --rm --entrypoint golangci-lint flat-search run ./...

tests:
	docker-compose run --rm --entrypoint go flat-search test ./...

code-coverage:
	@docker-compose run --rm --no-deps --entrypoint bash flat-search -xeuo pipefail -c "go test -coverprofile=/tmp/c.out ./...; \
	go tool cover -html=/tmp/c.out -o coverage.html; \
	"
	open coverage.html
