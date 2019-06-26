#define new line
define nl


endef

# This is how we want to name the binary output
BINARY=flat-search

# These are the values we want to pass for Version and BuildTime
VERSION=0.0.2
BUILD_TIME=`date +%FT%T%z`

.DEFAULT_GOAL := build

pre-install: export GO111MODULE=on
pre-install:
	go mod download
	go get -u -d github.com/golang-migrate/migrate/cli github.com/go-sql-driver/mysql
	go build -tags 'mysql' -o ${GOPATH}/bin/migrate github.com/golang-migrate/migrate/cli

build: export GO111MODULE=on
build:
	go mod verify
	gofmt -w ./
	go build -o ${BINARY} main.go

.PHONY: install
install: export GO111MODULE=on
install:
	go mod download
	go install ./...

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

migrate:
	read -p "Enter Database user: " USER; \
	read -p "Enter Database password: " PASSWORD; \
	read -p "Enter Database host: " HOST; \
	read -p "Enter Database port: " PORT; \
	read -p "Enter Database name: " NAME; \
	${GOPATH}/bin/migrate -database mysql://$$USER:$$PASSWORD@tcp\($$HOST:$$PORT\)/$$NAME -path ./migrations/ up

check:
	$(error "fix linter errors first:${LINTER_ERRORS}")

docker-compose:
	docker-compose down
	docker-compose -f docker-compose-migration.yml down
	docker-compose rm
	docker-compose up -d db
	docker-compose -f docker-compose-migration.yml build flat-migration
	docker-compose -f docker-compose-migration.yml up flat-migration
	docker-compose build flat-search
	docker-compose create flat-search
	docker-compose stop

docker-start:
	docker-compose start

docker-stop:
	docker-compose stop

