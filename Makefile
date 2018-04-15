#define new line
define nl


endef

# This is how we want to name the binary output
BINARY=flat-search

# These are the values we want to pass for Version and BuildTime
VERSION=0.0.2
BUILD_TIME=`date +%FT%T%z`

VENDOR_DIR=$(shell ls -d vendor | tail -n 1)
# check if vendor folder exists
ifneq (, $(VENDOR_DIR))
#check by linter
LINTER_ERRORS=$(subst vendor,$(nl)vendor,$(shell find ./ -path ./vendor -prune -o -type d,l -exec golint {} \;))
ifneq ("${LINTER_ERRORS}", "")
build: check
endif
endif

.DEFAULT_GOAL := build

pre-install:
	go get -u github.com/golang/dep/cmd/dep
	go get -u -d github.com/golang-migrate/migrate/cli github.com/go-sql-driver/mysql
	go build -tags 'mysql' -o ${GOPATH}/bin/migrate github.com/golang-migrate/migrate/cli
	${GOPATH}/bin/dep ensure

build:
	gofmt -w ./
	go build -o ${BINARY} main.go

.PHONY: install
install:
	${GOPATH}/bin/dep ensure
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

run.docker:
	docker-compose down
	docker-compose -f docker-compose-migration.yml down
	docker-compose rm
	docker-compose up -d db
	docker-compose -f docker-compose-migration.yml build flat-migration
	docker-compose -f docker-compose-migration.yml up flat-migration
	docker-compose build flat-search
	docker-compose create flat-search
	docker-compose stop

docker.start:
	docker-compose start

docker.stop:
	docker-compose stop

