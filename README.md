# Clodhopper

[![coverage report](https://gitlab.com/l3montree/crypto-koi/crypto-koi-api/badges/main/coverage.svg)](https://gitlab.com/l3montree/crypto-koi/crypto-koi-api/-/commits/main)

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

Example:

```sh
go run cmd/crypto-koi-cli/main.go [-drawPrimaryColor] <tokenId>
```

If the -drawPrimaryColor flag is provided, the image will contain the primary koi color in the top left corner. This color can be used by client side applications to modify the user interface colors accordingly.