.PHONY: update-proto deb_package_all
update-proto:
	protoc  --go_out=. --go-grpc_out=. common.proto repository.proto user.proto health.proto signal.proto
	mv github.com/EdmilsonRodrigues/ophelia-ci/* .
	rm -rf github.com
	./update_python_proto.bash

deb_package_all:
	cd interface && \
	 make build-deb
	cd client && \
	 make build && make package_deb
	cd server && \
	 make build && make package_deb
