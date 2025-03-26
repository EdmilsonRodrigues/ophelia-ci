FROM python:3.12-slim

WORKDIR /app

RUN pip install uv

COPY pyproject.toml uv.lock* ./

RUN uv pip install --system --no-cache-dir -r pyproject.toml

RUN mkdir -p ophelia_ci_interface/services 

COPY *.proto ophelia-ci-interface/services/

COPY interface/src/. .

ENV OPHELIA_CI_UNPACKED=True

ENTRYPOINT ["uvicorn"]

CMD ["ophelia_ci_interface.main:app", "--host", "0.0.0.0", "--port", "8000"]
EXPOSE 8000
