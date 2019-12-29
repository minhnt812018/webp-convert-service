FROM golang:latest

ENV libwebp_ver=libwebp-1.0.3

RUN apt-get -y update \
 && apt-get install -y git curl wget libjpeg-dev libpng-dev libtiff-dev libgif-dev libtool autoconf automake make gcc g++ imagemagick

WORKDIR /usr/local/webp
RUN wget https://storage.googleapis.com/downloads.webmproject.org/releases/webp/$libwebp_ver.tar.gz \
      && tar -xvzf $libwebp_ver.tar.gz \
      && cd $libwebp_ver \
      && ./configure --enable-everything \
      && make \
      && make install \
      && cd .. \
      && rm -rf $libwebp_ver

ENV PATH $PATH:/usr/local/webp/bin

RUN ldconfig

RUN mkdir -p /src
ADD src/main.go /src/main.go
WORKDIR /src

EXPOSE 80

CMD ["go","run","main.go"]
