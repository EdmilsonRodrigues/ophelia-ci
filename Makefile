.PHONY: build, package, build_and_package
build:
	go build -o deb-packaging/DEBIAN/usr/bin/ophelia-ci-server server/main.go
package:
	rm -rf dist/*
	./package.bash
	mv ophelia* dist/
build_and_package:
	make build && make package

