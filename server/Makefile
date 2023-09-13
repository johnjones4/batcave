PROJECT=$(shell basename $(shell cd .. && pwd))
ARCH=arm64
OS=linux
TAG=ghcr.io/johnjones4/${PROJECT}-${OS}-${ARCH}

postgres:
	docker run -e POSTGRES_PASSWORD=dev -e POSTGRES_USER=dev -e POSTGRES_DB=dev -p 5432:5432 ankane/pgvector

table:
	psql dev -h localhost -d dev < pgschema.sql

run:
	HTTP_HOST=:8080 DATABASE_URL=postgresql://dev:dev@localhost/dev SERVICE_CONFIG=./data/service_configuration.json go run . 

container:
	docker build -t ${TAG} --platform ${OS}/${ARCH} .
	docker push ${TAG}:latest
	docker image rm ${TAG}:latest
