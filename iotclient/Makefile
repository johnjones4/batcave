build:
	docker run --platform linux/arm/v7 --rm -v "${PWD}":/usr/src/app -w /usr/src/app -e CGO_ENABLED=1 -e GOOS=linux -e GOARCH=arm64 golang:1.21 go build -v -o bin/halclient .
