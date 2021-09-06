FROM golang:1.17 as build-stage

RUN apt-get update \
 && DEBIAN_FRONTEND=noninteractive \
    apt-get install --no-install-recommends --assume-yes \
      protobuf-compiler

RUN curl -sSL https://github.com/grpc/grpc-web/releases/download/1.2.1/\
protoc-gen-grpc-web-1.2.1-linux-x86_64 -o /usr/local/bin/protoc-gen-grpc-web && \
  chmod +x /usr/local/bin/protoc-gen-grpc-web

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

RUN export PATH="$PATH:$(go env GOPATH)/bin"

COPY go.mod /app/go.mod
COPY ./cmd /app/cmd/
COPY ./internal /app/internal/
COPY ./proto /app/proto/
COPY ./third_party /app/third_party/

WORKDIR /app
RUN go run ./cmd/main.go --all

FROM scratch as export-stage
COPY --from=build-stage /app/output /