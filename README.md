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

## Development

The api endpoints and models are generated based on the swagger 2.0 spec using the ["go-swagger"](https://goswagger.io) tool. After installing `go-swagger` run `swagger generate server -f ~$GOPATH/src/github.com/uhuchain/uhuchain-api/swagger/swagger.yaml -A uhuchain-api` from the root of this repository.

