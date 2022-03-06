FROM golang:1.17 AS build

WORKDIR /src/server/

COPY go.mod go.sum ./
RUN go mod download
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY ./cmd ./internal ./spec ./
COPY ./cmd/ ./cmd/
COPY ./internal/ ./internal/
COPY ./spec/ ./spec/

RUN apt-get update \
&& DEBIAN_FRONTEND=noninteractive \
    apt-get install --no-install-recommends --assume-yes \
      protobuf-compiler
    
RUN protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative spec/spec.proto

RUN go build -o /src/server/server ./cmd/server/routing_server.go

FROM ubuntu:latest

COPY --from=build /src/server/server /usr/local/bin/navserver
#COPY testRoute.txt /usr/local/bin/navserver/

RUN ["chmod", "+x", "/usr/local/bin/navserver"]