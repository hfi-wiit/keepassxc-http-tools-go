# vim: set ts=2 sw=2 noexpandtab :
.DELETE_ON_ERROR:

# build the binary in local build directory
.PHONY: build
build: build/kpht

build/kpht: main.go $(wildcard cmd/*.go)
	mkdir -p build
	go build -o $@ main.go

# remove local build directory
.PHONY: clean
clean:
	rm -rf build
