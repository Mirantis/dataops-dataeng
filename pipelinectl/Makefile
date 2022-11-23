GO_IMAGE := golang:1.16.5
DOCKER_FLAGS := -v $(MKFILE_DIR):/go/$(MOD_PATH) -w /go/$(MOD_PATH) --rm -t
GO_FLAGS     := $(DOCKER_FLAGS) -u $$(id -u):$$(id -g) -e GOCACHE=/tmp/go=cache -e GOOS=$(GOOS)
GO := docker run $(GO_FLAGS) $(GO_IMAGE) go

.PHONY: binary
binary:
	$(GO) build -o bin/pipelinectl $(MOD_PATH)

.PHONY: test 
test: 
	$(GO) test -v -cover $(MOD_PATH)/...