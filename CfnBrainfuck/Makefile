GO_TRANSFORM_PKG = $(wildcard transform/*.go)
GO_CLI_PKG = $(wildcard cli/*.go)
GO_LAMBDA_PKG = $(wildcard lambda/*.go)

NAME=CfnBrainfuck

.PHONY: all
all: ${NAME} ${NAME}.zip

.PHONY: lambda
lambda: ${NAME}.zip

${NAME}.zip: handler
	zip ${NAME}.zip handler

handler: $(GO_LAMBDA_PKG) $(GO_TRANSFORM_PKG)
	go build -o handler github.com/KablamoOSS/cfn-macros/${NAME}/lambda

.PHONY: cli
cli: ${NAME}

${NAME}: $(GO_CLI_PKG) $(GO_TRANSFORM_PKG)
	go build -o ${NAME} github.com/KablamoOSS/cfn-macros/${NAME}/cli

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm -f handler ${NAME} ${NAME}.zip
