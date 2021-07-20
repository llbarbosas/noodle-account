PROJECTNAME=$(shell basename "$(PWD)")
SHELL := /bin/bash
RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(RUN_ARGS):;@:)

default: help

## start: Start server in development mode.
start: cmd/start.go
	@go run cmd/start.go

## start-prod: Start server in production mode.
start-prod: cmd/start.go
	@ENV=prod go run cmd/start.go

## build: Build server.
build: cmd/start.go
	@go build -o ./server cmd/start.go

## docker-build: Build and run docker image.
docker-build: Dockerfile
	@docker build -t noodle-account . && docker run -p 3001:3001 noodle-account

## test: Run test suites
test: 
	@go test ./...

## todo: Show all to-do comment locations
todo:
	@grep -rnw '.' -e '// TODO:' --exclude=Makefile    

## local-ssl: Generate local ssl certificate
local-ssl: 
	@openssl req -x509 -out https/certs/ssl.cert -keyout https/certs/ssl.key -newkey rsa:2048 -nodes -sha256 -subj '/CN=localhost' -extensions EXT -config <( printf "[dn]\nCN=localhost\n[req]\ndistinguished_name = dn\n[EXT]\nsubjectAltName=DNS:localhost\nkeyUsage=digitalSignature\nextendedKeyUsage=serverAuth")

.PHONY: help
help: Makefile
	@echo "Available commands on "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'