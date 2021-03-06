# CryptoKoi Api

[![coverage report](https://gitlab.com/l3montree/crypto-koi/crypto-koi-api/badges/main/coverage.svg)](https://gitlab.com/l3montree/crypto-koi/crypto-koi-api/-/commits/main)

# Addresses

## Local Development

Owner: `0x2bb6335AC37c468c626D18C9915A8Cc7c36D76e7` (managed by timbastin)

Mumbai Testnet Smart-Contract Address: `0x5dF55eA9e0035d588F755d92f18bD207b1435bDc`

## Dev Environment
Owner: `0x2bb6335AC37c468c626D18C9915A8Cc7c36D76e7` (managed by timbastin)

Mumbai Testnet Smart-Contract Address: `0x2f158340c811c37284748fD2FFE298ebEB2F8c0e`


## Prod Environment
Owner:

Polygon Smart-Contract Address:

## Requirements

1. Docker
2. Geth (https://geth.ethereum.org/docs/install-and-build/installing-geth)
3. gqlgen (`go install github.com/99designs/gqlgen`)
## Setup

Install gqlgen - this is the utility to generate the go structs based on the defined graphql types.

```sh
go install github.com/99designs/gqlgen
```

Generate a private ECDSA key using the following command:

```sh
openssl ecparam -name prime256v1 -genkey -noout -out ./key.pem && openssl ec -in ./key.pem -pubout -out ./public.pem
```

Rename the `.env.example` to `.env`

```sh
cp .env.example .env
```

Check the provided example values in the environment file. The environment variables do include absolute paths. Therefore some changes need to be made.

Start everything using `make`.

```sh
make
```

This will start the docker containers, regenerates the types and start the server.

## Web3

The web3 integration is build using typescript and etherjs. For local testing hardhat is used. To start a local blockchain (hardhat) use the `make web3` command. This will start a blockchain network accessible at `http://localhost:8545`. To deploy the smart contract `CryptoKoi` onto the chain, the `make deploy` command can be used. This will first:

1. Compile the contract
2. Deploy it

## CLI Usage

The application ships with a cli to generate kois using a token id. Make sure to set the `BASE_IMAGE_PATH` environment variable to the absolute path to the folder `./images/raw`.


### Draw the image of a specific token

Draw the image of a specific token. The token needs to be either in `HEX` or `DECIMAL` format (e.g "1238af213hhffff", "12356524234234").
Example:

```sh
go run cmd/crypto-koi-cli/main.go [-drawPrimaryColor] [-debug] draw <tokenId>
```

If the -drawPrimaryColor flag is provided, the image will contain the primary koi color in the top left corner. This color can be used by client side applications to modify the user interface colors accordingly.

### Register a random user

This can be helpful when testing the different client side interface colors.

```sh
go run cmd/crypto-koi-cli/main.go [-amount] [-debug] register <tokenId>
```


The amount of the users to be registered can be provided using the amount flag.
