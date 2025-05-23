.PHONY: all clean install

BINARY_NAME := gweather
GO_BIN_DIR := $(shell go env GOPATH)/bin

all: $(BINARY_NAME)

$(BINARY_NAME):
	go build -o $(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)

install:
	go install

uninstall:
	rm -f $(GO_BIN_DIR)/$(BINARY_NAME)
