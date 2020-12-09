.PHONY: lint

lint:
	@golangci-lint run --timeout=5m \
		--enable=bodyclose \
		--enable=dogsled \
		--enable=dupl \
		--enable=godox \
		--enable=gomnd \
		--enable=gosec \
		--enable=gocritic \
		--enable=goerr113 \
		--enable=gocyclo \
		--enable=gocognit \
		--enable=gofmt \
		--enable=maligned \
		--enable=unparam \
		--enable=prealloc \
		--enable=wsl \
		--enable=testpackage \
		--enable=exportloopref \
		--enable=exhaustive \
		--enable=sqlclosecheck \
		--out-format=colored-line-number
