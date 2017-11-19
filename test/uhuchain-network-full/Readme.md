# Hyperledger fabric example

First try out of a Hyperledger Fabric network based on the 'build your first network' tutorial. This sample includes the use of CAs, Apache Kafka / Zookeeper, and multiple orderer peers. 

For reference check the official tutorial: ["Build Your First Network"](http://hyperledger-fabric.readthedocs.io/en/latest/build_network.html)

## Prequisites

Make sure you got the platform-specific binaries installed just as descibed ["here"](https://hyperledger-fabric.readthedocs.io/en/latest/samples.html)

## Layout

The network should be used by two (famous) organizations - Cyberdyne Cooperation and Umbrella Cooperation.

- Zookeeper 1
- Zookeeper 2
- Zookeeper 3
- Kafka 1
- Kafka 2
- Kafka 3
- Kafka 4
- CA Insurance A
- CA Insurance B
- CA Insurance C
- Orderer
    - orderer.insurancea.uhuchain.com
    - orderer.insuranceb.uhuchain.com
    - orderer.insurancec.uhuchain.com
- Peers
    - Insurance A
        - peer0.insurancea.uhuchain.com
        - peer1.insurancea.uhuchain.com
    - Insurance B
        - peer0.insurancea.uhuchain.com
        - peer1.insurancea.uhuchain.com
    - Insurance C
        - peer0.insurancea.uhuchain.com
        - peer1.insurancea.uhuchain.com

## Generate the certificates
All crypto materials are generated using the `cryptogen` tool from the platform-specific binaries.

Create a configuration file `crypto-config.yml` containing the orderer and peer nodes with their specific domain names and run the tool.

`../bin/cryptogen generate --config=./crypto-config.yaml`

It generates a folder `crypto-config` containing all the cryptographic materials.

## Generate the transaction artifacts

Next setup the configuration file `configtx.yaml` for creating the transaction artifacts.

- genesis block
- channel artifact
- anchor peers

The parameter `MSPDir` links the defined organisation to the certificates that have been generated in the previous step.

The parameter `ID` sets the unique ID of the MSP. Do **not** use any special character for it (_,-,+ etc.).

In the orderer section set the `OrdererType` to kafka and add at least two kafka brokers with their hostname and port (Kafka brokers will be defined in the docker compose file).

The profile section defines the configurtion for generating the different artifacts.

Export the base directory of the project into the variable `FABRIC_CFG_PATH`, since it is required by the `configtxgen` tool.

`export FABRIC_CFG_PATH=$PWD`
`mkdir channel-artifacts`

### genesis block
Create the genesis block using the `ThreeOrgsOrdererGenesis` profile.

`../bin/configtxgen -profile ThreeOrgsOrdererGenesis -outputBlock ./channel-artifacts/genesis.block`

### channel
Create the channel between the parties that have been defined in a profile `ThreeOrgsChannel` in `configtx.yaml`. Channel name needs to be set by exporting the `CHANNEL_NAME` variable and must satisfy the following requirements:
1. Contain only lower case ASCII alphanumerics, dots '.' and dashes '-'
2. Are shorter than 250 characters.
3. Start with a letter

`export CHANNEL_NAME=car-ledger`
`../bin/configtxgen -profile ThreeOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME`

### anchor peers
After creating the channel, at least one anchorpeer needs to be setup for each organisation. An anchorpeer can always be seen by all other participants of the network. The `-asOrg` parameter should corelated to the organization name used in the `configtx.yaml` (`Organizations` -> `<Organisation>` -> `Name`).

Create anchorpeer for the first organization (Insurance A here)
`../bin/configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/InsuranceAMSPanchors.tx -channelID $CHANNEL_NAME -asOrg InsuranceA`

Create anchorpeer for the second organization (Insurance B here)

`../bin/configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/InsuranceBMSPanchors.tx -channelID $CHANNEL_NAME -asOrg InsuranceB`

Create anchorpeer for the second organization (Insurance C here)

`../bin/configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/InsuranceCMSPanchors.tx -channelID $CHANNEL_NAME -asOrg InsuranceC`

## Setup docker compose
The file `docker-compose-senexum.yaml` contains the definitions for all containerized components for the blockchain network. Enviroment settings and mounted volumes need changed to be changed to match the files (certificates and keys) and configurations created in the previous steps. 


### Bring up the network
Before bringing up the network, make sure that in the `peer_base.yaml` the parameter `CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=<network>` matches the network name defined in `docker-compose-senexum.yaml` in the `networks` section. Note that docker compose prefixes the network name. To check for network names use `docker network ls` while the composition is running.

Example here: `CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=uhuchainnetwork_uhuchain`

Start the network using `docker-compose` command.

`CHANNEL_NAME=$CHANNEL_NAME TIMEOUT=1000 docker-compose -f docker-compose-full.yaml up -d`

Bring down the network using `docker-compose -f docker-compose-full.yaml down`.

Use `docker rm $(docker ps -a -q)` to remove all containers including the ones running the chaincode.

## Running the script

For testing purposes, the docker compose file also contains an addition container `cli` with includes several tools to interact with the network from the console.

It also includes an end-2-end test script from the `byfn` tutorial.

Within the test Script we need to change domain names, path to certificates, and MSP IDs (also the one in the endorsment policy!)

Enter the CLI container uding `docker exec -it cli bash` and execute the end-2-end test with `./scripts/script.sh`.

See the logs of the CLI container by using `docker logs -f cli`.

### Manual execution

**Invoke**
`peer chaincode invoke -o orderer.insurancea.uhuchain.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/orderer.uhuchain.com/orderers/orderer.insurancea.uhuchain.com/msp/tlscacerts/tlsca.orderer.uhuchain.com-cert.pem -C car-ledger -n mycc -c '{"Args":["invoke","a","b","100"]}'`

**Query**
`peer chaincode query -C car-ledger -n mycc -c '{"Args":["query","b"]}'`