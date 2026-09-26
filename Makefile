.PHONY: test
test:
	go test -v -covermode=count -coverprofile=coverage.out ./...

.PHONY: cover
cover: test
	go tool cover -html=coverage.out

.PHONY: fmt
fmt:
	gofumpt -w .

.PHONY: check
check:
	golangci-lint run
	@echo
	gofumpt -l .

.PHONY: clean
clean:
	rm -rf bin dist coverage.out

.PHONY: help
help:
	@echo "Usage: make <target>"
	@echo
	@echo "Targets:"
	@echo "  test   Run tests with coverage (coverage.out)"
	@echo "  cover  Open the HTML coverage report"
	@echo "  fmt    Format with gofumpt        (go install mvdan.cc/gofumpt@latest)"
	@echo "  check  Run golangci-lint + gofumpt (go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest)"
	@echo "  clean  Remove bin/, dist/ and coverage.out"
	@echo "  help   Show this help"
