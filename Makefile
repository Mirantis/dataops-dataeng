MKFILE_DIR := $(abspath $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST))))))
IMAGE_NAME ?= mirantiseng/dataeng
IMAGE_TAG ?= latest

.PHONY: container 
container:
	docker build $(MKFILE_DIR) -t $(IMAGE_NAME):$(IMAGE_TAG)