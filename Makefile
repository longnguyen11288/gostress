ARCH=$(shell uname -m)
PATH=bin/$(ARCH)

build: build_client build_server

build_client:
	go build -o $(PATH)/client pool.go client.go

build_server:
	go build -o $(PATH)/server pool.go server.go
