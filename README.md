# JHU Deliverbot Server

This server provides admin control of the robots, as well as assigning them tasks and providing routing information for navigation.

## Installation & Running

- Install Go
- Clone this repository

To get OSRM mapping data, run either `osrm_data_pipeline.ps1` in Powershell or `osrm_data_pipeline.sh` in Bash.  This should create a `/data` folder under `/OSRM` and populate it with OSRM map data files.

Make a copy of the config.template.json (`cp config.template.json config.json`). See [Certificates](#certificates) for how to generate certificates for testing &mdash; if you use the local testing example you won't need to modify the config at all.

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

## Certificates

> **Note** <br>
If you just need certificates for running the simulator and server locally on your computer, see [local testing](#generating-certs-for-local-testing). Certificates *must* be generated in the `certs/` folder for them to be accessible in Docker.

We rely on mTLS both to encryt communications between the server and the robots, as well as for authentication. All security is at the transport level: if you aren't authorized to use the server, you simply cannot connect at all. There are no passwords or other tokens - having an appropriate certificate provides proof of identity.

certgen.py is a tool for generating mTLS certificates

Setting up mTLS security involves two steps
 - generating a root CA (certificate authority) which will issue all the mTLS certificates.
 - using the CA to issue certificates to the server and robots

Each certificate (including the root CA) has an accompaning private key, which must be kept secret. Since the C++ and Go libaries for gRPC don't support password encrypted private keys, neither do their respective standard libraries, the private keys are kept in plain text. *This means the `.key` files must be kept safe!*

> **Warning** <br>
> DO NOT check certificates or keys into Git - if you do, you *must* revoke the certificate, or the entire chain. If you don't know how to do this, ask someone!

> **Note** <br>
>  Certificates *must* be generated in the `certs/` folder for them to be accessible in Docker.

### Examples

To generating a root CA named "deliverbot_ca" (if you just need to add a new robot, don't do this, just use the existing CA to issue a new certificate):
```
python .\certgen.py -p <ca password> ca -n deliverbot_ca
```

To issue a cert for a hypothetical server named "deliverbot_server" running at `fangornbane.com`:
```
python .\certgen.py -p <cert password> issue -n deliverbot_server -c deliverbot_ca --ca_password <ca password> -u fangornbane.com
```

The certgen utility also allows changing the expiration period (default is 30 days) and modifying the subject organization info for the certificates.

The CA certificate needs to be set as the root CA on both the server and the robots. The CA key file and password should be kept secret, they are needed to issue additional certificates.

### Generating certs for local testing:
```
python .\certgen.py create_ca -n local_test_ca
python .\certgen.py issue -n local_test_server -c local_test_ca -u localhost
python .\certgen.py issue -n local_test_robot -c local_test_ca
```
