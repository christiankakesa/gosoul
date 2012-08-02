# GoSoul client Makefile
#
# Christian Kakesa <christian.kakesa@gmail.com>

all: fmt gosoul_client

gosoul_client: bin/gosoul_client.go
	go build -o bin/$@ $<

fmt:
	gofmt -w *.go
	gofmt -w bin/*.go

clean:
	rm -rf gosoul_client
