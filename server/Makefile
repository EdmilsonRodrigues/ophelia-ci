TAG ?= latest

.PHONY: build, package_deb, package_snap, package_docker, build_and_package, changelog
build:
	go build -o deb-packaging/usr/bin/ophelia-ci-server .
package_deb:
	./package_deb.bash
package_snap:
	./package_snap.bash
package_docker:
	docker build -t edmilsonrodrigues/ophelia-ci-server:$(TAG) -f Dockerfile ..
build_and_package:
	rm -rf dist/ophelia-ci-server*
	make build && make package_deb && make package_snap
	mv ophelia* dist/
changelog:
	./changelog.bash
