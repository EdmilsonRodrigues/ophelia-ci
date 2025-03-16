FROM golang:1.24 AS builder

WORKDIR /app

COPY . /app

RUN make build

FROM ubuntu AS debug

COPY --from=builder /app/dist/*.deb /app/

# RUN dpkg -i /app/*.deb
ENTRYPOINT ["bash"]