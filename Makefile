# GoSoul client Makefile
#
# Christian Kakesa <christian.kakesa@gmail.com>

.PHONY: all fmt clean gosoul

all: fmt gosoul

gosoul: gosoul.go
	go build -o $@ $<

fmt:
	gofmt -w **/*.go

clean:
	rm -rf gosoul *.gz
