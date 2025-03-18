.PHONY: lint mypy run build
lint:
	ruff format
	ruff check . --fix
	ruff format

mypy:
	mypy .

run:
	uvicorn app.main:app --reload

build:
	briefcase create
	briefcase build
	briefcase package
	briefcase package --target fedora:40
	briefcase package linux flatpak
