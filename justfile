set tempdir := "./tmp"
set quiet := true
set dotenv-load := true

branch := shell('git rev-parse --abbrev-ref HEAD')
tag := shell('git describe --tags --abbrev=0 $(git rev-list --tags --max-count=1 --date-order main)')
version := if branch == "main" { tag } else { tag+ "-dev" }

# Clean artifacts, update dependencies, generate code, lint, run tests, and compile the application
default: clean dep generate lint test compile

# Clean build artifacts and temporary files
clean:
    rm -f ./dist/*
    rm -rf ./static/assets/*.css ./static/assets/*.js ./static/assets/*.woff ./static/assets/*.woff2 ./static/*.html ./static/*.ico ./static/*.png

# Ensure API and UI dependencies are up to date and valid
dep:
    go mod tidy
    go mod verify
    cd ./ui && pnpm install --frozen-lockfile

# Run database migrations and generate API code (OAPI and SQLC)
generate: gen-queries gen-api

# Run database migrations
migrations:
    echo "Running migrations..."
    go run . database migrate

# Generate SQL code from SQL files using sqlc
gen-queries:
    echo "Generating queries..."
    go tool sqlc generate --file ./etc/db/sqlc.yaml


# Generate API code from OpenAPI spec using oapi-codegen
gen-api:
    echo "Generating API..."
    go tool oapi-codegen -config ./etc/api/server.yaml ./etc/api/spec.jsonc

# Build the UI using pnpm
gen-ui:
    echo "Building UI...";
    cd ./ui && pnpm run build


# Run linters for both API and UI code
lint: api-lint ui-lint

# Check UI code for issues with eslint
ui-lint:
    cd ./ui && pnpm run lint

# Check API for code issues with golangci-lint
api-lint:
    golangci-lint run

# Run API and UI unit tests
test: api-test ui-test

# Run UI tests using pnpm and vitest
ui-test:
    cd ./ui && pnpm run test

# Run Go tests with coverage and output results in JSON format for tparse
api-test:
    go test ./... -cover -json -count 1 | tparse -trimpath github.com/jsnfwlr/filamate/

coverage: api-coverage

# Generate an API code coverage report and save it as an HTML file in the static directory
api-coverage:
    go test ./... -covermode=count -coverprofile=coverage.out -json -count=1 | tparse -pass -trimpath github.com/jsnfwlr/filamate/ || rm coverage.out
    go tool cover -html=coverage.out -o static/coverage.html && sed -i 's|rgb(128, 128, 128)|#00CCFF|; s|rgb(116, 140, 131)|#00D2F9|; s|rgb(104, 152, 134)|#00D8F3|; s|rgb(92, 164, 137)|#00DEEE|; s|rgb(80, 176, 140)|#00E3E8|; s|rgb(68, 188, 143)|#00E9E3|; s|rgb(56, 200, 146)|#00EEDD|; s|rgb(44, 212, 149)|#00F4D8|; s|rgb(32, 224, 152)|#00F9D2|; s|rgb(20, 236, 155)|#00FFCC|;' static/coverage.html
    rm coverage.out
    echo "Coverage report can be accessed at /coverage.html if the server is running."

# Build the application binary with the embedded UI assets
compile:
    cd ./ui && pnpm run build
    go build -ldflags "-X github.com/jsnfwlr/filamate/internal/cmd.Version={{version}}" -o ./dist/filamate .

# Run the application with live reloading for both API and UI using Air and pnpm
run:
    cd ./ui; pnpm run build-live &
    air -c .air.toml &

# Watch for changes in the UI and rebuild it automatically
ui-run:
	cd ./ui; pnpm run build-live

# Run the application with live reloading using Air
api-run:
    air -c .air.toml

# Install latest versions of development tools: air, tparse, sqlc, and oapi-codegen
tools:
    echo "Installing development tools..."
    echo "- air"
    go install github.com/air-verse/air@latest
    echo "- tparse"
    go install github.com/mfridman/tparse@latest
    echo "- sqlc"
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
    echo "- oapi-codegen"
    go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

# Build the Docker image for the application
build:
	docker build -t filamate:latest .

# Start the application and its dependencies using Docker Compose
up:
	docker compose -f ./compose.yaml up -d

# Stop and remove the application containers, networks, and volumes created by Docker Compose
down:
	docker compose -f ./compose.yaml down
