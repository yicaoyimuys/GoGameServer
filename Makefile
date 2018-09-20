.PHONY: .FORCE
GO=go

SRC_DIR = ./src

NEW_GOPATH = $(GOPATH):$(shell pwd)
GOPATH := $(NEW_GOPATH)

all:
	$(GO) install servive/connector

clean:
	rm -rf bin pkg release
	rm -rf logs/*

fmt:
	$(GO) fmt $(SRC_DIR)/...

vendor_init:
	cd $(SRC_DIR) && govendor init

vendor_addExternal:
	cd $(SRC_DIR) && govendor add +external

publish_linux:
	GOOS=linux GOARCH=amd64 $(GO) build -o release/connector servive/connector
	
publish_windows:
	GOOS=windows GOARCH=amd64 $(GO) build -o release/connector.exe connector
	
publish_mac:
	GOOS=darwin GOARCH=amd64 $(GO) build -o release/connector connector
	