.PHONY: build, package_deb, build_and_package
build:
	go build -o deb-packaging/usr/bin/ophelia-ci-server server/main.go
package_deb:
	./package_deb.bash
package_snap:
	./package_snap.bash
build_and_package:
	rm -rf dist/*
	make build && make package_deb && make package_snap
	mv ophelia* dist/

