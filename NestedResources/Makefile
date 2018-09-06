GO_TRANSFORM_PKG = $(wildcard transform/*.go)
GO_CLI_PKG = $(wildcard cli/*.go)
GO_LAMBDA_PKG = $(wildcard lambda/*.go)

.PHONY: all
all: nested-resources nested-resources.zip

nested-resources.zip: handler
	zip nested-resources.zip handler

handler: $(GO_LAMBDA_PKG) $(GO_TRANSFORM_PKG)
	go build -o handler github.com/KablamoOSS/nested-resources/lambda 

nested-resources: $(GO_CLI_PKG) $(GO_TRANSFORM_PKG)
	go build -o nested-resources github.com/KablamoOSS/nested-resources/cli 

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm -f handler nested-resources nested-resources.zip
