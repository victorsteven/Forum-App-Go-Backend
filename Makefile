GOFILES = $(shell find . -name '*.go')

default: build

workdir:
	mkdir -p workdir

build: workdir/contacts

build-native: $(GOFILES)
	go build -o workdir/native-contacts .

workdir/contacts: $(GOFILES)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o workdir/contacts .
