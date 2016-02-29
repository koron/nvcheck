default: test

test:
	go test ./...

lint:
	@echo ""
	go vet ./...
	@echo ""
	golint ./...

cyclo:
	@echo ""
	-gocyclo -over 9 -avg .

report:
	@echo ""
	@echo "misspell"
	@find . -name "*.go" | xargs misspell
	@echo ""
	gocyclo -over 14 .
	@echo ""
	go vet ./...
	@echo ""
	golint ./...

.PHONY: default test lint report
