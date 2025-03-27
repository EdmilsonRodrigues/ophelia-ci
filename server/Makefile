TAG ?= latest

.PHONY: build, deb, snap, docker, rock

build:
	go build -o deb-packaging/usr/bin/ophelia-ci-server .

deb:
	./package_deb.bash

snap:
	./package_snap.bash

rock:
	rockcraft pack
	mv *.rock ../dist/

docker:
	docker build -t edmilsonrodrigues/ophelia-ci-server:$(TAG) -f Dockerfile ..

