migrate-up:
	docker run -v "${CURDIR}/migrations:/migrations" --network host migrate/migrate -path=/migrations/ -database "mysql://root:test100500@tcp(localhost:3306)/flat-search" up

migrate-down:
	docker run -v "${CURDIR}/migrations:/migrations" --network host migrate/migrate -path=/migrations/ -database "mysql://root:test100500@tcp(localhost:3306)/flat-search" down -all
