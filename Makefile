
gosoul_client: gosoul_client.go
	go build $^

fmt:
	gofmt -w *.go
	gofmt -w gosoul/*.go

clean:
	rm -rf gosoul_client
