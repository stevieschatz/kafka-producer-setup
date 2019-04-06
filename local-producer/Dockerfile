FROM golang:alpine as builder
RUN apk update && apk add curl git pkgconfig alpine-sdk bash
WORKDIR /root
RUN git clone https://github.com/edenhill/librdkafka.git
WORKDIR /root/librdkafka
RUN /root/librdkafka/configure
RUN make
RUN make install
#For golang applications
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN mkdir /build
ADD . /build/
WORKDIR /build
ENV PKG_CONFIG_PATH=$PKG_CONFIG_PATH:/usr/lib/pkgconfig/
RUN CGO_ENABLED=1 GOOS=linux go build -tags static_all -o local-producer .
CMD ["./local-producer"]
