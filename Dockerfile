FROM golang:1.18.5 AS build

RUN mkdir -p /opt/eavesdropper

WORKDIR /opt/eavesdropper

COPY go.mod go.sum Makefile /opt/eavesdropper/
RUN make mod

COPY . /opt/eavesdropper/

RUN make build-eavesdropper

EXPOSE 8080

CMD ["./eavesdropper", "--config", "config.yaml"]
