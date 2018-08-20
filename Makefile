# Makefile

GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_BUILD_RACE=$(GO_CMD) build -race
GO_TEST=$(GO_CMD) test
GO_TEST_VERBOSE=$(GO_CMD) test -v
GO_TEST_COVER=$(GO_CMD) test -cover
GO_INSTALL=$(GO_CMD) install -v

SERVER_BIN=crawler
SERVER_DIR=.
SERVER_MAIN=main.go

SOURCE_PKG_DIR=.
SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

all: test build-server run

build-server:
	@echo "==> Building server ...";
	@$(GO_INSTALL)
	@$(GO_BUILD)
	@chmod 755 $(SERVER_BIN)

test:
	@echo "==> Running tests ...";
	@$(GO_TEST_COVER) $(SOURCE_PKG_DIR)

run:
	./$(SERVER_BIN)
