# GoSoul client Makefile
#
# Christian Kakesa <christian.kakesa@gmail.com>

gosoul_client: gosoul_client.go fmt
	go build $<

fmt:
	gofmt -w *.go
	gofmt -w gosoul/*.go

clean:
	rm -rf gosoul_client
