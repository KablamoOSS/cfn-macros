GO_TRANSFORM_PKG = $(wildcard transform/*.go)
GO_CLI_PKG = $(wildcard cli/*.go)
GO_LAMBDA_PKG = $(wildcard lambda/*.go)

NAME=LocalTransform

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
