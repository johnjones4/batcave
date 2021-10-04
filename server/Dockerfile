FROM golang:1.16
WORKDIR /usr/src/api
COPY . .
RUN go get .
RUN go build .
VOLUME [ "/config" ]
CMD ["/usr/src/api/main"]