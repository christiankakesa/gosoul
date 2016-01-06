# GoSoul client Makefile
#
# Christian Kakesa <christian.kakesa@gmail.com>

all: fmt gosoul

gosoul: gosoul.go
	go build -o $@ $<

fmt:
	gofmt -w *.go
	gofmt -w **/*.go

clean:
	rm -rf gosoul *.gz
