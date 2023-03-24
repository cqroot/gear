.PHONY: install
install:
	go install .

.PHONY: test
test:
	go test -v -covermode=count -coverprofile=coverage.out ./...

.PHONY: cover
cover: test
	go tool cover -html=coverage.out

.PHONY: check
check:
	@echo '******************************'
	golangci-lint run
	@echo
	@echo '******************************'
	gofumpt -l .
