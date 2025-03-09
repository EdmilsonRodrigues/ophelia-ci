.PHONY: build, package, build_and_package
build:
	go build -o deb-packaging/DEBIAN/usr/bin/ophelia-ci-server main.go
package:
	./package.bash
build_and_package:
	make build && make package

