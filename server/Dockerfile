FROM golang:1.17

WORKDIR /usr/src/app

COPY . .

WORKDIR /usr/src/app/api
RUN go get .
RUN go build .


VOLUME [ "/data" ]


CMD ["/usr/src/app/api/main"]
