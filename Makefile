# Copyright (C) 2026 Keith Chu <cqroot@outlook.com>
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.

.PHONY: test
test:
	go test -v -covermode count -coverprofile coverage.out ./...

.PHONY: cover
cover: test
	go tool cover -html coverage.out

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
