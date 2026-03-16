BINARY := omachy
VERSION ?= dev
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

.PHONY: build run clean tidy test test-setup

build:
	go build $(LDFLAGS) -o $(BINARY) .

run: build
	./$(BINARY) install

clean:
	rm -f $(BINARY)

tidy:
	go mod tidy

test:
	go test ./...

test-setup:
	./test/integration/setup-vm.sh

test-integration:
	./test/integration/run-tests.sh
