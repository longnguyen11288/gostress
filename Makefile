build: build_client build_server

build_client:
	go build -o client pool.go client.go

build_server:
	go build -o server pool.go server.go