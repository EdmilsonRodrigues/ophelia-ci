.PHONY: update-proto
update-proto:
	protoc  --go_out=. --go-grpc_out=. repository.proto
	mv github.com/EdmilsonRodrigues/ophelia-ci/* .
	rm -rf github.com
	cd interface/src/ophelia_ci_interface/services \
	 && python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. repository.proto \
	 && sed -i "s/import repository_pb2 as repository__pb2/import ophelia_ci_interface.services.repository_pb2 as repository__pb2/" repository_pb2_grpc.py \
	 && sed -i "s/DESCRIPTOR.*/&  # noqa: E501/" repository_pb2.py
	
