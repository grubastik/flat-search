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
LINTER_ERRORS=$(subst vendor,$(nl)vendor,$(shell find ./vendor/github.com/grubastik/flat-search/ -type d,l -exec golint {} \;))
ifneq ("${LINTER_ERRORS}", "")
build: check
endif
endif

.DEFAULT_GOAL := build

pre-install:
	curl https://glide.sh/get | sh
	go get -u github.com/mattes/migrate

build:
	gofmt -w ./vendor/github.com/grubastik/flat-search/
	go build -o ${BINARY} main.go

.PHONY: install
install:
	${GOPATH}/bin/glide install
	go install ./...

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

migrate:
	${GOPATH}/bin/migrate

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

