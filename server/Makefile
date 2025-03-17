.PHONY: build, package_deb, build_and_package
build:
	go build -o deb-packaging/usr/bin/ophelia-ci-server .
package_deb:
	./package_deb.bash
package_snap:
	./package_snap.bash
build_and_package:
	rm -rf dist/ophelia-ci-server*
	make build && make package_deb && make package_snap
	mv ophelia* dist/

