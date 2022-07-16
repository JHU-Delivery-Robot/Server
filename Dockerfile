FROM golang:1.18 AS build

WORKDIR /src/server/

COPY go.mod go.sum ./
RUN go mod download
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY ./cmd/ ./cmd/
COPY ./internal/ ./internal/
COPY ./protocols/ ./protocols/

RUN apt-get update \
&& DEBIAN_FRONTEND=noninteractive \
    apt-get install --no-install-recommends --assume-yes \
      protobuf-compiler

RUN protoc --proto_path="protocols" --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative protocols/routing.proto
RUN protoc --proto_path="protocols" --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative protocols/development.proto

RUN go build -o /src/server/navserver ./cmd/server/

FROM ubuntu:latest

COPY --from=build /src/server/navserver /usr/local/bin/navserver/

WORKDIR /usr/local/bin/navserver
ENTRYPOINT [ "./navserver" ]
