GO     ?= go
GOFMT  ?= gofmt -s
GOPATH ?= $($(GO) env GOPATH)

BIN  ?= /usr/local/bin
UDEV ?= /etc/udev/rules.d

DASH   ?= $(if $(GOOS),_,)
GOOS   ?=
GOARCH ?=

EXEC    := xavierSrv$(DASH)$(GOOS)$(DASH)$(GOARCH)
GOOUT   := $(shell pwd)/bin
GOFILES := main.go

.PHONY: help
help:
	@echo "Make Routines:"
	@echo " - \"\"                print Make routines list"
	@echo " - build             creates the binary, Cross-Compiling ex : 'make build GOOS='linux' GOARCH='arm''"
	@echo " - release           Build for Linux 64/32/arm and OpenBSD 64/32"
	@echo " - clean             delete build files"
	@echo " - install           install the binary and service"
	@echo " - remove            remove binary and service"
	@echo " - upgrade           build and update installed binary"

.PHONY: install
install: build
	sudo cp $(GOOUT)/$(EXEC) $(BIN)/$(EXEC)
	# Add service
	make clean

.PHONY: remove
remove:
	sudo rm -rf $(BIN)/$(EXEC)
	# Remove service

.PHONY: upgrade
upgrade: build
	sudo cp $(GOOUT)/$(EXEC) $(BIN)/$(EXEC)
	# Restart service
	make clean

.PHONY: fmt
fmt:
	$(GOFMT) -e -w $(GOFILES)

.PHONY: godep
godep:
	$(GO) get github.com/BurntSushi/toml
	$(GO) get github.com/xhit/go-simple-mail

.PHONY: build
build: godep
	GOARCH=$(GOARCH) GOOS=$(GOOS) $(GO) build -o $(GOOUT)/$(EXEC) -a

.PHONY: buildall
buildall:
	GOARCH=amd64 GOOS=linux make build
	GOARCH=amd64 GOOS=openbsd make build
	GOARCH=arm GOOS=linux make build
	GOARCH=arm GOOS=openbsd make build
	GOARCH=386 GOOS=linux make build

.PHONY: clean
clean:
	$(GO) clean
	rm -rf $(GOOUT)/*
