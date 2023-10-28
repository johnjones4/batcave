FROM debian:bookworm

RUN apt-get update && apt-get install -y \
  wget \
  libasound2-dev \
  sox \
  libwebkit2gtk-4.0-dev \
  libgtk-3-dev

RUN wget https://go.dev/dl/go1.21.3.linux-arm64.tar.gz
RUN tar -C /usr/local -xzf go1.21.3.linux-arm64.tar.gz
