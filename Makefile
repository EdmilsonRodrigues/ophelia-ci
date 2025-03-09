.PHONY: build
build:
	go build -o packaging/DEBIAN/usr/bin/ophelia-ci-server main.go
