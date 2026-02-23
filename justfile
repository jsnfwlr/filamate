default: clean gomod generate lint test build

clean:
	@rm -f ../bin/filamate

gomod:
	@go mod tidy
	@go mod verify

generate: db api ui

db:
	@echo "Running migrations..."
	@go run . database migrate
	@echo "Generating queries..."
	@go tool sqlc generate --file ./etc/db/sqlc.yaml

api:
	@echo "Generating API..."
	@go tool oapi-codegen -config ./etc/api/server.yaml ./etc/api/spec.jsonc

ui:
	@echo "Building UI..."
	@cd ./ui && npm install && npm run build
	@touch main.go

lint:
    @golangci-lint run

test:
    @go test ./... -cover -json -count 1 | tparse -trimpath git.phalacee.com/jsnfwlr/filamate/

build:
	@go build -o ../bin/filamate .

run:
	API_PORT="9767" POSTGRES_PORT="5444" air -c .air.toml

coverage:
	@go test ./... -covermode=count -coverprofile=coverage.out -json -count=1 | tparse -pass -trimpath github.com/jsnfwlr/filamate/ || rm coverage.out
	@go tool cover -html=coverage.out -o static/coverage.html && sed -i 's|rgb(128, 128, 128)|#00CCFF|; s|rgb(116, 140, 131)|#00D2F9|; s|rgb(104, 152, 134)|#00D8F3|; s|rgb(92, 164, 137)|#00DEEE|; s|rgb(80, 176, 140)|#00E3E8|; s|rgb(68, 188, 143)|#00E9E3|; s|rgb(56, 200, 146)|#00EEDD|; s|rgb(44, 212, 149)|#00F4D8|; s|rgb(32, 224, 152)|#00F9D2|; s|rgb(20, 236, 155)|#00FFCC|;' static/coverage.html
	@rm coverage.out
	@echo "Coverage report can be accessed at /coverage.html if the server is running."

tools:
	@echo "Installing development tools..."
	@echo "- tparse"
	@go install github.com/mfridman/tparse@latest
	@echo "- sqlc"
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@echo "- oapi-codegen"
	@go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
