FROM golang:1.16

RUN apt-get update && apt-get install -y \
  wget \
  build-essential \
  cmake

WORKDIR /usr/src

RUN wget https://github.com/mit-nlp/MITIE/archive/refs/tags/v0.7.tar.gz
RUN tar zxvf v0.7.tar.gz

WORKDIR /usr/src/MITIE-0.7/mitielib
RUN make
RUN make install
RUN ldconfig


WORKDIR /usr/src/MITIE-0.7
RUN make MITIE-models


WORKDIR /usr/src/api

COPY . .

WORKDIR /usr/src/api/api
RUN go get .
RUN go build .


VOLUME [ "/data" ]


CMD ["/usr/src/api/api/main"]
