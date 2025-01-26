API_GATEWAY_BINARY=apiGatewayBinary
AUTH_BINARY=authBinary
LOGGER_BINARY=loggerBinary
MAP_STORAGE_BINARY=mapStorageBinary

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker compose up -d
	@echo "Docker images started!"

## up_build: stops docker compose (if running), builds all projects and starts docker compose
up_build: build_api_gateway build_auth build_map_storage build_logger
	@echo "Stopping docker images (if running...)"
	docker compose down
	@echo "Building (when required) and starting docker images..."
	docker compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker compose down
	@echo "Done!"

## build_api_gateway: builds the api_gateway binary as a linux executable
build_api_gateway:
	@echo "Building API gateway binary..."
	cd ./api-gateway-service && mkdir -p bin && env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o bin/${API_GATEWAY_BINARY} ./cmd/api
	@echo "Done!"

## build_auth: builds the auth binary as a linux executable
build_auth:
	@echo "Building auth binary..."
	cd ./authentication-service && mkdir -p bin && env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o bin/${AUTH_BINARY} ./cmd/api
	@echo "Done!"

## build_logger: builds the logger binary as a linux executable
build_logger:
	@echo "Building logger binary..."
	cd ./logger-service && mkdir -p bin && env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o bin/${LOGGER_BINARY} ./cmd/api
	@echo "Done!"

## build_map_storage: builds the map_storage binary as a linux executable
build_map_storage:
	@echo "Building map storage binary..."
	cd ./map-storage-service && mkdir -p bin && env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o bin/${MAP_STORAGE_BINARY} ./cmd/api
	@echo "Done!"

## build_logger: builds the logger binary as a linux executable
build_logger:
	@echo "Building logger binary..."
	cd ./logger-service && mkdir -p bin && env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o bin/${LOGGER_BINARY} ./cmd/api
	@echo "Done!"







