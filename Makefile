BINARY = conduit

.PHONY: build
build:
	dep ensure	
	make test
	go build -o ${BINARY}

.PHONY: test
test:
	go test ./...

.PHONY: run
run:
	make build
	./${BINARY}

.PHONY: release
release:
	glide install
	GOOS=linux GOARCH=amd64 go build -o build/${BINARY}-linux-amd64 .;
	GOOS=darwin GOARCH=amd64 go build -o build/${BINARY}-darwin-amd64 .;
	GOOS=windows GOARCH=amd64 go build -o build/${BINARY}-windows-amd64.exe .;
