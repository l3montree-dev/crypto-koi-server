MAKEFLAGS += -j2

.PHONY: run docker codegen deploy


run: docker codegen
	go run cmd/crypto-koi-api/main.go

docker: 
	docker-compose up -d

codegen: 
	gqlgen generate

node_modules:
	cd web3 && npm i && cd ..

deploy:
	cd web3 && npm run deploy && cd ..

start-web3:
	cd web3 && npm run start && cd ..