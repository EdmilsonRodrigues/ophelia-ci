.PHONY: update-proto deb_package_all
update-proto:
	protoc  --go_out=. --go-grpc_out=. repository.proto
	protoc  --go_out=. --go-grpc_out=. user.proto
	protoc  --go_out=. --go-grpc_out=. health.proto
	mv github.com/EdmilsonRodrigues/ophelia-ci/* .
	rm -rf github.com
	cd interface/src/ophelia_ci_interface/services \
	 && python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. repository.proto \
	 && python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. user.proto \
	 && python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. health.proto \
	 && sed -i "s/import repository_pb2 as repository__pb2/import ophelia_ci_interface.services.repository_pb2 as repository__pb2/" repository_pb2_grpc.py \
	 && sed -i "s/DESCRIPTOR.*/&  # noqa: E501/" repository_pb2.py \
	 && sed -i "s/DESCRIPTOR.*/&  # noqa: E501/" user_pb2.py

deb_package_all:
	cd interface && \
	 make build-deb
	cd client && \
	 make build && make package_deb
	cd server && \
	 make build && make package_deb
