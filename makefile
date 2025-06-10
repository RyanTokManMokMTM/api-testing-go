gen-workflow:
	go run utils/tool/cmd/main.go

test/api-junit:
	ginkgo --fail-fast --junit-report report.xml --output-dir ./reports run ./test

test/api-json:
	ginkgo --fail-fast --json-report report.json --output-dir ./reports run ./test

# Linting commands
lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

lint-install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Test coverage commands
test-coverage:
	go test -coverprofile=coverage.out -covermode=atomic ./utils/... ./api_test_runner

test-coverage-clean:
	rm -f coverage.out coverage.html

# Development helpers
dev-setup: lint-install
	go mod tidy
	go mod download

