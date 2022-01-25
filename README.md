# Clodhopper

[![coverage report](https://gitlab.com/l3montree/cryptogotchi/clodhopper/badges/main/coverage.svg)](https://gitlab.com/l3montree/cryptogotchi/clodhopper/-/commits/main)

## Setup

Generate a private ECDSA key using the following command:

```sh
openssl ecparam -name prime256v1 -genkey -noout -out ./testdata/key.pem && openssl ec -in ./testdata/key.pem -pubout -out ./testdata/public.pem
```