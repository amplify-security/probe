
.PHONY: mod
mod:
	go mod tidy

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: clean
clean:
	rm -rf vendor

.PHONY: test
test:
	go test -race ./...