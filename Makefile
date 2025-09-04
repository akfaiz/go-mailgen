COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html
MIN_COVERAGE  := 90

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
.PHONY: lint
lint:
	@echo "Linting code..."
	@golangci-lint run --timeout 5m

# Install dependencies and tools
.PHONY: install
install:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v -cover -race ./... -coverprofile=$(COVERAGE_FILE) -coverpkg=./...

.PHONY: coverage
coverage:
	@echo "Generating test coverage report..."
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@go tool cover -func=$(COVERAGE_FILE) | tee coverage.txt
	@echo "Coverage HTML report generated: $(COVERAGE_HTML)"
	@open $(COVERAGE_HTML)

.PHONY: coverage-check
coverage-check:
	@COVERAGE=$$(go tool cover -func=$(COVERAGE_FILE) | grep total | awk '{print $$3}' | sed 's/%//'); \
	RESULT=$$(echo "$$COVERAGE < $(MIN_COVERAGE)" | bc); \
	if [ "$$RESULT" -eq 1 ]; then \
		echo "Coverage is below $(MIN_COVERAGE)%: $$COVERAGE%"; \
		exit 1; \
	else \
		echo "Coverage is sufficient: $$COVERAGE%"; \
	fi

.PHONY: update-examples
update-examples:
	@echo "Updating example files..."
	go run ./examples/main.go
