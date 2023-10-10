FROM golang:1.21 as builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o hal9000

FROM alpine:3.10
WORKDIR /app
RUN mkdir /data
ENV SERVICE_CONFIG=/data/service_configuration.json
COPY start.sh start.sh
COPY --from=builder /app/hal9000 /app/hal9000
RUN apk add --no-cache curl
ENTRYPOINT [ "/app/start.sh" ]
