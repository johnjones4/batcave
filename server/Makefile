PROJECT=$(shell basename $(shell cd .. && pwd))
ARCH=amd64
OS=linux
TAG=ghcr.io/johnjones4/${PROJECT}-${OS}-${ARCH}

postgres:
	docker run -e POSTGRES_PASSWORD=dev -e POSTGRES_USER=dev -e POSTGRES_DB=dev -p 5432:5432 ankane/pgvector

table:
	psql dev -h localhost -d dev < pgschema.sql

seed:
	psql dev -h localhost -d dev < seed_intent.sql
	psql dev -h localhost -d dev < seed_users.sql

run:
	HTTP_HOST=:8080 DATABASE_URL=postgresql://dev:dev@localhost/dev SERVICE_CONFIG=./data/service_configuration.json go run . 

container:
	go test -cover ./...
	docker build -t ${TAG} --platform ${OS}/${ARCH} .
	docker push ${TAG}:latest
	docker image rm ${TAG}:latest

intent-seed-openai:
	python3 scripts/generate_openai_intent_seed.py > seed_intent.sql

.PHONEY:

mocks: .PHONEY
	~/go/bin/mockgen -source=./core/intent_types.go -package=mocks  > ./mocks/intent_types.go
	~/go/bin/mockgen -source=./core/message_types.go -package=mocks  > ./mocks/message_types.go
	~/go/bin/mockgen -source=./core/service_types.go -package=mocks  > ./mocks/service_types.go
	~/go/bin/mockgen -source=./core/service_types.go -package=mocks  > ./mocks/service_types.go
	~/go/bin/mockgen -source=./core/store_types.go -package=mocks  > ./mocks/store_types.go
	~/go/bin/mockgen -source=./core/types.go -package=mocks  > ./mocks/types.go
	~/go/bin/mockgen -source=${HOME}/go/pkg/mod/github.com/sirupsen/logrus@v1.9.3/logrus.go -package=mocks  > ./mocks/logrus.go

test-coverage:
	go test -coverprofile cover.out ./...
	go tool cover -html=cover.out
