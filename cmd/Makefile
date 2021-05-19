default: build

.PHONY: build
build:
	@go build -v -mod=vendor -o doraemon

.PHONY: install
install:
	@go build -v -mod=vendor -o ${GOPATH}/bin/doraemon