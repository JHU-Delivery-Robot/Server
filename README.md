# JHU Deliverbot Server

This server provides admin control of the robots, as well as assigning them tasks and providing routing information for navigation.

## Installation & Running

- Install Go
- Clone this repository

To get OSRM mapping data, run either `osrm_data_pipeline.ps1` in Powershell or `osrm_data_pipeline.sh` in Bash.  This should create a `/data` folder under `/OSRM` and populate it with OSRM map data files.

To run this server, run:
- `docker compose build`
- `docker compose up`

To force re-building the image, remove the containers using `docker compose down` and then remove the images with `docker rmi <image ids>` where `<image ids>` are relevant image IDs, which can be found using `docker images`.

## Building gRPC

We use protobufs and gRPC as the server protocol. When the gRPC service definitions are changed, the Go stubs needed to be re-generated. This can be done using the commands
```
protoc --proto_path="protocols" --go_out="protocols" --go_opt=paths=source_relative --go-grpc_out="protocols" --go_opt= --go-grpc_opt=paths=source_relative protocols/routing.proto
protoc --proto_path="protocols" --go_out="protocols" --go_opt=paths=source_relative --go-grpc_out="protocols" --go_opt= --go-grpc_opt=paths=source_relative protocols/development.proto
```
