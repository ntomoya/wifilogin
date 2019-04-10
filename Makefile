VERSION := 0.1
REVISION := $(shell git rev-parse --short HEAD)

SRCS := $(shell find . -type f -name '*.go')
LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

all: build

build: $(SRCS)
	go install -a -tags netgo -installsuffix netgo $(LDFLAGS)

run: build
	bin/$(NAME)
