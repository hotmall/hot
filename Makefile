PACKAGES = $(shell go list ./... | grep -v '/vendor/')

PACKAGE = github.com/hotmall/hot
COMMIT_HASH = $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE = $(shell date +%FT%T%z)
SERVICE_VERSION = $(shell cat VERSION)
GO_VERSION = $(shell go version)

ldflagsversion = -X $(PACKAGE)/commands.CommitHash=$(COMMIT_HASH) -X $(PACKAGE)/commands.BuildDate=$(BUILD_DATE) -X $(PACKAGE)/commands.Version=$(SERVICE_VERSION) -X "$(PACKAGE)/commands.GoVersion=$(GO_VERSION)" -s -w

all: install

install:	
	go install -v -ldflags '$(ldflagsversion)'

test:
	go test $(PACKAGES)
