# Makefile for calculated values
.DEFAULT_GOAL    := cbp

# CBP_ROOT       :=$(shell git rev-parse --show-toplevel)
VERSION          := $(shell git describe --tags --dirty)
COMMIT_HASH      := $(shell git rev-parse --short HEAD 2>/dev/null)
DATE             := $(shell date "+%Y-%m-%d")
# BUILD_PLATFORM := $(shell uname -a | awk '{print tolower($1);}')

BUILD_INFO := -X main.version=${VERSION} -X main.commitHash=${COMMIT_HASH} -X main.buildDate=${DATE}

SOURCES := main.go version.go doc.go Makefile


cbp: $(SOURCES)
	go build -ldflags "${BUILD_INFO}" -o cbp .

minimal: $(SOURCES)
	@# create a static image, from:
	@#  https://github.com/docker-library/golang/issues/152
	go build -ldflags "-d -s -w ${BUILD_INFO}" -tags netgo -o cbp .
.PHONY: minimal

doc.go: main.go mkdoc.sh
	go generate -v -x .
