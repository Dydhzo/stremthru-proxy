NAME = stremthru-proxy

.PHONY: all clean fmt test build run docker-build docker-run

all: build

clean:
	rm -rf $(NAME)

fmt:
	go fmt ./...

test:
	STREMTHRU_ENV=test go test -v ./...

build: clean
	go build

run:
	go run .

docker-build:
	docker build -t stremthru-proxy:latest .

docker-run:
	docker run --rm -it --name $(NAME) \
		-p 8080:8080 \
		stremthru-proxy:latest
