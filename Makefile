.PHONY: clean
clearn:
	@go vet ./...
	@go fmt ./...

# brew install golang-migrate
.PHONY: migrateadd
migrateadd:
	@migrate create -seq -ext=.sql -dir=./migrations $(tablename)

GREENLENGHT_DB_DSN=postgres://greenlight:pa55word@127.0.0.1:5432/greenlight?sslmode=disable
.PHONY: migraterun
migraterun:
	@migrate -path=./migrations -database=$(GREENLENGHT_DB_DSN) up

# make migratedown version=1
.PHONY: migratedown
migratedown:
	@migrate -path=./migrations -database=$(GREENLENGHT_DB_DSN) down $(version)

.PHONY: migratedownall
migratedownall:
	@migrate -path=./migrations -database=$(GREENLENGHT_DB_DSN) down 

.PHONY: run
run:
	@go run ./cmd/api/...

.PHONY: api-test
api-test: 
	@hurl api.hurl
