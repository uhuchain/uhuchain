# Uhuchain API server

The uhuchain api server offers a RESTful API to interact with the uhuchain ledger.

## Requirements

* Go 1.9.x
* Docker
* Docker Compose

## Installation

Use the `make` file to install all required dependencies.

Note: Run the `make` command always from the root of the project.

```
make depend-install
```

## Run 

The repo includes a hyperledger test network consisting of three parties and an ordering service.

Running `make integration-test` will 
* start the uhuchain test network using `docker-compose` 
* create an `uhuchain-server` container exposing port 3333
* run the integration tests in `test/integration`
* compile and install the `uhuchain-server` binary
* start the `uhuchain-server` on port 3333

You can access the API through `http://localhost:3333/v1/status`

Example: `curl 'http://localhost:3333/v1/status'`

### Behind the scenes
The `uhuchain-server` binary has several command line options. Use `uhuchain-server --help` for details.

```
Usage:
  uhuchain-server [OPTIONS]

Uhuchain simple REST API for cars

Application Options:
      --scheme=            the listeners to enable, this can be repeated and defaults to the schemes in the swagger spec
      --cleanup-timeout=   grace period for which to wait before shutting down the server (default: 10s)
      --max-header-size=   controls the maximum number of bytes the server will read parsing the request header's keys and values, including the request line. It does not limit the
                           size of the request body. (default: 1MiB)
      --socket-path=       the unix socket to listen on (default: /var/run/uhuchain.sock)
      --host=              the IP to listen on (default: localhost) [$HOST]
      --port=              the port to listen on for insecure connections, defaults to a random value [$PORT]
      --listen-limit=      limit the number of outstanding requests
      --keep-alive=        sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download) (default: 3m)
      --read-timeout=      maximum duration before timing out read of the request (default: 30s)
      --write-timeout=     maximum duration before timing out write of the response (default: 60s)
      --tls-host=          the IP to listen on for tls, when not specified it's the same as --host [$TLS_HOST]
      --tls-port=          the port to listen on for secure connections, defaults to a random value [$TLS_PORT]
      --tls-certificate=   the certificate to use for secure connections [$TLS_CERTIFICATE]
      --tls-key=           the private key to use for secure conections [$TLS_PRIVATE_KEY]
      --tls-ca=            the certificate authority file to be used with mutual tls auth [$TLS_CA_CERTIFICATE]
      --tls-listen-limit=  limit the number of outstanding requests
      --tls-keep-alive=    sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download)
      --tls-read-timeout=  maximum duration before timing out read of the request
      --tls-write-timeout= maximum duration before timing out write of the response

Help Options:
  -h, --help               Show this help message
```
Example: `uhuchain-server --scheme=http --host=0.0.0.0 --port=3333`

## Development

The api endpoints and models are generated based on the swagger 2.0 spec using the ["go-swagger"](https://goswagger.io) tool. After installing `go-swagger` run `swagger generate server -f $GOPATH/src/github.com/uhuchain/uhuchain-api/swagger/swagger.yaml -A uhuchain-api` from the root of this repository.

The generator will replace **all** files in the following directories:

* `cmd`
* `models`
* `restapi`

except for

* `restapi/configure_uhuchain.go`
* `restapi/handler/*`

These files should be used for the actual implmentation of the endpoints.

## Testing

For development and testing, login to the `uhuchain-server` container with `docker exec -it uhuchain-server bash`.

Run the integration test manually with `go test -v ./test/integration/`.

### Updating chaincode

The `prepare.sh` scripts supports upgrading the chaincode to a new version by running `./test/uhuchain-network-dev/scripts/prepare.sh car-ledger 1 upgrade 1.1`. 

TODO: Change the default path for the chaincode to the uhuchain chaincode. Currently the default fabric samples are used.

### docker compose setup

For development and testing a docker compose network is defined under `test/uhuchain-network-dev`. Please see the `Readme.md` there for more details.

