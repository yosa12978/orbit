.DEFAULT_GOAL := build

build:
	@go mod tidy
	@go build -o ./bin/orbit .

run: build
	@ORBIT_ADDR=0.0.0.0:3000 ./bin/orbit
