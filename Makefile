.PHONY: clean build

BINARY_NAME=crypto-http
BUILD := `git rev-parse HEAD`
go=GO111MODULE=on go
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

build: darwin  darwin-arm64 linux64

darwin-arm64:
	@GOARCH=arm64 GOOS=darwin ${go} build  ${LDFLAGS} -x -o bin/${BINARY_NAME}-arm64 ./cmd/main.go

darwin:
	@GOARCH=amd64 GOOS=darwin ${go} build  ${LDFLAGS} -x -o bin/${BINARY_NAME}-darwin ./cmd/main.go

linux64:
	@GOARCH=amd64 GOOS=linux ${go} build ${LDFLAGS} -x -o bin/${BINARY_NAME}-linux64 ./cmd/main.go

format:
	@gofmt -l -w ./pkg/crypto ./cmd

run_test:
	@${go} test ./pkg/crypto -v

clean:
	go clean
	rm ./bin/${BINARY_NAME}-*