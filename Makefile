.PHONY: build, package, build_and_package
build:
	go build -o packaging/DEBIAN/usr/bin/ophelia-ci-server main.go
package:
	./package.bash
build_and_package:
	make build && make package

