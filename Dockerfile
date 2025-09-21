FROM golang:1.24 AS builder

WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY core ./core
COPY internal ./internal
COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./stremthru-proxy -a -ldflags '-extldflags "-static"'

FROM alpine

WORKDIR /app

COPY --from=builder /workspace/stremthru-proxy ./stremthru-proxy

EXPOSE 8080

ENTRYPOINT ["./stremthru-proxy"]
