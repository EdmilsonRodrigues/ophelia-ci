.PHONY: lint mypy run build build-deb
lint:
	ruff format
	ruff check . --fix
	ruff format

mypy:
	mypy .

run:
	cd src && \
	OPHELIA_CI_DEBUG=True uvicorn ophelia_ci_interface.main:app --reload

build:
	briefcase create
	briefcase build
	briefcase package
	briefcase package --target fedora:40
	briefcase package linux flatpak
	mv dist/* ../dist

build-deb:
	briefcase create
	briefcase build
	briefcase package
	mv dist/* ../dist
