## server:                     run the server
server:
	@go build -o build/bin/server && ./build/bin/server
## docker build:               build docker image
docker-build:
	@docker build -t envoy-authz:build .
## docker run:                 run docker image
docker-run:
	@make docker-build
	@docker run --rm --name envoy-authz -p 9001:9001 --network local -e GRPC_SERVER_ADDR=:9001 envoy-authz:build
