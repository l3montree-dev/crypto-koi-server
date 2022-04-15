MAKEFLAGS += -j2

.PHONY: run docker codegen deploy

all: docker codegen
	go run cmd/crypto-koi-api/main.go

run::
	go run cmd/crypto-koi-api/main.go

test: codegen
	go test ./...

docker: 
	docker-compose up -d

codegen: graphql abi

abi: contracts
	abigen --bin=./contracts_CryptoKoi_sol_CryptoKoi.bin --abi=./contracts_CryptoKoi_sol_CryptoKoi.abi --pkg=cryptokoi --type=CryptoKoiBinding --out=internal/cryptokoi/cryptokoi_binding.go

contracts: contracts_CryptoKoi_sol_CryptoKoi.bin contracts_CryptoKoi_sol_CryptoKoi.abi
	rm -f @*

contracts_CryptoKoi_sol_CryptoKoi.bin: contracts/CryptoKoi.sol node_modules
	npx solc --include-path node_modules/ --base-path . --bin contracts/CryptoKoi.sol

contracts_CryptoKoi_sol_CryptoKoi.abi: contracts/CryptoKoi.sol node_modules
	npx solc --include-path node_modules/ --base-path . --abi contracts/CryptoKoi.sol

graphql: graph/schema.graphqls
	gqlgen generate

deploy: node_modules
	npm run deploy

blockchain: node_modules
	npm run start

node_modules: package.json
	npm install