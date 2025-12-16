FROM golang:1.24-alpine AS builder 
WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .

RUN mkdir -p ./bin
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o ./bin/fotosouk ./cmd/main.go

FROM gcr.io/distroless/static-debian12:nonroot 
WORKDIR /app 
COPY --from=builder /app/bin ./bin
EXPOSE 3000
USER nonroot:nonroot
CMD [ "./bin/fotosouk" ]