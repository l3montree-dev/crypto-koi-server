# Clodhopper

[![coverage report](https://gitlab.com/l3montree/cryptogotchi/clodhopper/badges/main/coverage.svg)](https://gitlab.com/l3montree/cryptogotchi/clodhopper/-/commits/main)

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