BINARY := omachy
VERSION ?= dev
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

.PHONY: build run clean tidy

build:
	go build $(LDFLAGS) -o $(BINARY) .

run: build
	./$(BINARY) install

clean:
	rm -f $(BINARY)

tidy:
	go mod tidy
