FROM golang:1.17 as builder

WORKDIR /opt/

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
# RUN go test -cover ./...

ENV GOOS linux
ARG GOARCH
ENV GOARCH ${GOARCH:-amd64}
ENV CGO_ENABLED=0

ARG VERSION
RUN go build -v -o yaae cmd/yaae/main.go

FROM alpine:latest

EXPOSE 9100
ENTRYPOINT ["yaae"]
RUN addgroup -g 1000 exporter && \
    adduser -u 1000 -D -G exporter exporter -h /exporter

WORKDIR /exporter/


RUN apk --no-cache add ca-certificates
COPY --from=builder /opt/yaae /usr/local/bin/yaae
USER exporter
