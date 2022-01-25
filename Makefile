MAKEFLAGS += -j2

.PHONY: run docker codegen


run: docker codegen
	go run cmd/clodhopper-server/main.go

docker: 
	docker-compose up -d

codegen: 
	gqlgen generate