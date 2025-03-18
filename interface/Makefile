.PHONY: lint mypy run
lint:
	ruff format
	ruff check . --fix
	ruff format

mypy:
	mypy .

run:
	uvicorn app.main:app --reload
