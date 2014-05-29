GO_CMD=go
GOLINT_CMD=golint
GO_TEST=$(GO_CMD) test -v ./...
GO_VET=$(GO_CMD) vet ./...
GO_LINT=$(GOLINT_CMD) .

all:
	$(GO_VET)
	$(GO_LINT)
	$(GO_TEST)
