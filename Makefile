default: test

test:
	go test ./...

lint:
	@echo ""
	go vet ./...
	@echo ""
	golint ./...

report:
	@echo ""
	@echo "misspell"
	@find . -name "*.go" | xargs misspell
	@echo ""
	-gocyclo -over 9 -avg .
	@echo ""
	go vet ./...
	@echo ""
	golint ./...

.PHONY: default test lint report
