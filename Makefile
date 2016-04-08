GO_FILES := $(shell find . -type f -name '*.go' -not -path "./Godeps/*" -not -path "./vendor/*")
GO_PACKAGES := $(shell go list ./... | sed "s/github.com\/dmathieu\/gitest/./" | grep -v "^./vendor/")

build:
	go build -v $(GO_PACKAGES)

travis: tidy test

test: build
	go fmt $(GO_PACKAGES)
	go test -race -i $(GO_PACKAGES)
	go test -race -v $(GO_PACKAGES)

tidy: goimports
	test -z "$$(goimports -l -d $(GO_FILES) | tee /dev/stderr)"
	go vet $(GO_PACKAGES)

goimports:
	go get golang.org/x/tools/cmd/goimports
