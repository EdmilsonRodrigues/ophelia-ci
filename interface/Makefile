TAG ?= latest

.PHONY: lint mypy run build build-deb
lint:
	ruff format
	ruff check . --fix
	ruff format

mypy:
	mypy

run:
	cd src && \
	OPHELIA_CI_UNPACKED=True uvicorn ophelia_ci_interface.main:app --reload

docker:
	docker build -t edmilsonrodrigues/ophelia-ci-interface:$(TAG) -f Dockerfile ..

snap:
	./create_requirements.bash
	./update_versions.bash
	cd src && \
	snapcraft

rock:
	make mypy
	./create_requirements.bash
	./update_versions.bash
	cd src && \
	ROCKCRAFT_ENABLE_EXPERIMENTAL_EXTENSIONS=true rockcraft pack
	mv src/*.rock ../dist

build:
	./fix_pyproject.bash
	briefcase create
	briefcase build
	briefcase package
	briefcase package --target fedora:40
	briefcase package linux flatpak
	mv dist/* ../dist

build-deb:
	make mypy
	./fix_pyproject.bash
	briefcase create
	briefcase build
	briefcase package
	mv dist/* ../dist
