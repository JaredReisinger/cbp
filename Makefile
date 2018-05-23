# Makefile for calculated values
.DEFAULT_GOAL    := cbp

# CBP_ROOT       :=$(shell git rev-parse --show-toplevel)
VERSION          := $(shell git describe --tags --dirty)
COMMIT_HASH      := $(shell git rev-parse --short HEAD 2>/dev/null)
DATE             := $(shell date "+%Y-%m-%d")
# BUILD_PLATFORM := $(shell uname -a | awk '{print tolower($1);}')

GO_BUILD_LDFLAGS := -s -w -X main.version=${VERSION} -X main.commitHash=${COMMIT_HASH} -X main.buildDate=${DATE}

SOURCES := main.go version.go doc.go Makefile


cbp: $(SOURCES)
	go build -ldflags "${GO_BUILD_LDFLAGS}" -o cbp .

for-docker: $(SOURCES)
	go build -ldflags "-d ${GO_BUILD_LDFLAGS}" -tags netgo -o cbp .
.PHONY: for-docker

doc.go: main.go mkdoc.sh
	go generate -v -x .
