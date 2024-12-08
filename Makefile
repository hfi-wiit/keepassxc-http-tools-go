# vim: set ts=2 sw=2 noexpandtab :
# --contains gives short commit hash, while otherwise if there was a tag before
# we get something like v0.0.1-2-g2a674c0 (seems -2 = commits after v0.0.1)
VERSION ?= $(shell git describe --tags --always HEAD)

.DELETE_ON_ERROR:

# build the binary in local build directory
.PHONY: build
build: build/kpht

build/kpht: main.go $(wildcard cmd/*.go)
	mkdir -p build
# go tool nm build/kpht | grep Version | grep cmd
	go build -ldflags="-X 'keepassxc-http-tools-go/cmd.Version=$(VERSION)'" -o $@ main.go

# remove local build directory
.PHONY: clean
clean:
	rm -rf build
