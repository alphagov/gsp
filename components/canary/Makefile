CANARY_IMAGE_NAME = govsvc/gsp-canary
BUILD_TIMESTAMP = $(shell date +%s)

.PHONY: all
all: canary

canary:
	docker build -t $(CANARY_IMAGE_NAME) .

