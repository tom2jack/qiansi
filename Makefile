export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE := auto
LDFLAGS := -s -w

all: clean fmt package

swag:
	swag init

test:
	go test -v --cover ./...

clean:
	rm -f ./bin/*

fmt:
	go fmt ./...

package:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o ./bin/qiansi_server_linux_amd64

