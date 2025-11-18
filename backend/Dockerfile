# STOCKS/backend/Dockerfile
########################
# Builder stage (Debian-based so cgo + libsqlite3 compile works)
########################
FROM golang:1.21-bullseye AS builder
ENV CGO_ENABLED=1
WORKDIR /app

# install system deps needed to build go-sqlite3
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    libc6-dev \
    libsqlite3-dev \
    ca-certificates \
  && rm -rf /var/lib/apt/lists/*

# copy go modules first to leverage build cache
COPY go.mod go.sum ./
RUN go mod download

# copy sources and build
COPY . .

# build server binary (adjust main.go path if your main is in another file)
RUN go build -v -o /app/stocks-server ./main.go

########################
# Runtime stage (slim)
########################
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
  && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/stocks-server /usr/local/bin/stocks-server

RUN mkdir -p /data && chown -R 1000:1000 /data

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/stocks-server"]
