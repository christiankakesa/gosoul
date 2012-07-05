include $(GOROOT)/src/Make.$(GOARCH)

gosoul_client: fmt gosoul.$O gosoul_client.$O
	$(LD) -o gosoul_client gosoul_client.$O

fmt:
	gofmt -w *.go
	gofmt -w gosoul/*.go

gosoul_client.$O: gosoul_client.go
	$(GC) -I . -I gosoul/ gosoul_client.go

gosoul.$O: gosoul/*.go  
	$(GC) -I . -I gosoul/ gosoul/*.go

clean:
	rm -rf gosoul/*.$O *.$O gosoul_client
