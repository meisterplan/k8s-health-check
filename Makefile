build:
	go build -o check -v

test: build
	go test -v ./...

clean:
	go clean
	rm -f check

docker-build:
	rm -f check && docker run --rm -v "${PWD}":/usr/src/check -e GOCACHE=/tmp -u `id -u`:`id -g` -w /usr/src/check golang:1.10-alpine go build -o check -v

