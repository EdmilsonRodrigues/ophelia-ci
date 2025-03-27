TAG ?= latest

.PHONY: build, deb, snap, docker, rock

build:
	go build -o deb-packaging/usr/bin/ophelia-ci .

deb:
	./package_deb.bash

snap:
	./package_snap.bash

docker:
	docker build -t edmilsonrodrigues/ophelia-ci-client:$(TAG) -f Dockerfile ..
