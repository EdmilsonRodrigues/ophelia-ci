#!/bin/bash

PROTOS=("common.proto" "repository.proto" "user.proto" "health.proto" "signal.proto")

source .venv/bin/activate
cd interface/src/ophelia_ci_interface/services
python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. ${PROTOS[@]}
protoc --mypy_out=. ${PROTOS[@]}

for file in $(ls *.py); do
    sed -i "s/import common_pb2 as common__pb2/import ophelia_ci_interface.services.common_pb2 as common__pb2/" $file
    sed -i "s/import grpc/import grpc  # type: ignore[import-untyped]/" $file
    sed -i "s/    from grpc._utilities import first_version_is_lower/    from grpc._utilities import first_version_is_lower  # type: ignore[import-untyped]/" $file
    file_no_py=$(echo $file | sed 's/_grpc\.py//')
    sed -i "s/import $file_no_py as \(.*\)/import ophelia_ci_interface.services.$file_no_py as \1/" $file
done

