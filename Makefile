MAKEFLAGS += -j2

.PHONY: run docker codegen deploy


run: docker codegen
	go run cmd/crypto-koi-api/main.go

docker: 
	docker-compose up -d

codegen: 
	gqlgen generate

node_modules:
	npm i

deploy:
	npm run deploy

blockchain:
	npm run start