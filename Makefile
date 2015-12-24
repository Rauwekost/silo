#Makefile for Douane

#make sure gopath is there
ifndef GOPATH
$(warning You need to set up a GOPATH.  See the README file.)
endif

#settings
PROJECT := github.com/rauwekost/silo
PROJECT_DIR := $(shell go list -e -f '{{.Dir}}' $(PROJECT))
EXE				:= silo
USER			:= rauwekost
REPO			:= silo
NAME 			:= silo

UNIX_EXE := \
	darwin/amd64/$(EXE) \
	freebsd/amd64/$(EXE) \
	linux/amd64/$(EXE)
WIN_EXE := \
	windows/amd64/$(EXE).exe

#compresed state
COMPRESSED_EXE=$(UNIX_EXE:%=%.tar.bz2) $(WIN_EXE:%.exe=%.zip)
COMPRESSED_EXE_TARGETS=$(COMPRESSED_EXE:%=bin/%)

#version number & description
VERSION=$(shell git describe --abbrev=0 | egrep -o '([0-9]+\.){1,10}[0-9]+')
COMMIT=$(shell git rev-parse --short HEAD)
RELEASE_VERSION=${VERSION}-${COMMIT}

#upload command
UPLOAD_CMD = github-release upload --user $(USER) --repo $(REPO) --tag v$(RELEASE_VERSION) --name $(subst /,-,$(FILE)) --file bin/$(FILE)

#arm
bin/linux/arm/5/$(EXE):
	GOARM=5 GOARCH=arm GOOS=linux go build -o "$@" -ldflags "-X main.version=${RELEASE_VERSION}"
bin/linux/arm/7/$(EXE):
	GOARM=7 GOARCH=arm GOOS=linux go build -o "$@" -ldflags "-X main.version=${RELEASE_VERSION}"

#386
bin/darwin/386/$(EXE):
	GOARCH=386 GOOS=darwin go build -o "$@" -ldflags "-X main.version=${RELEASE_VERSION}"
bin/linux/386/$(EXE):
	GOARCH=386 GOOS=linux go build -o "$@" -ldflags "-X main.version=${RELEASE_VERSION}"
bin/windows/386/$(EXE):
	GOARCH=386 GOOS=windows go build -o "$@" -ldflags "-X main.version=${RELEASE_VERSION}"

#amd64
bin/freebsd/amd64/$(EXE):
	GOARCH=amd64 GOOS=freebsd go build -o "$@" -ldflags "-X main.version=${RELEASE_VERSION}"
bin/darwin/amd64/$(EXE):
	GOARCH=amd64 GOOS=darwin go build -o "$@" -ldflags "-X main.version=${RELEASE_VERSION}"
bin/linux/amd64/$(EXE):
	GOARCH=amd64 GOOS=linux go build -o "$@" -ldflags "-X main.version=${RELEASE_VERSION}"
bin/windows/amd64/$(EXE).exe:
	GOARCH=amd64 GOOS=windows go build -o "$@" -ldflags "-X main.version=${RELEASE_VERSION}"

# compressed artifacts, makes a huge difference (Go executable is ~9MB,
# after compressing ~2MB)
%.tar.bz2: %
	tar -jcvf "$<.tar.bz2" "$<"
%.zip: %.exe
	zip "$@" "$<"

bin/tmp/$(EXE):
	go build -o "$@" -ldflags "-X main.version=${RELEASE_VERSION}"

#phony
#------------------------------------------------------------------------------#
default: build

build:
	go build $(PROJECT)/...

format:
	gofmt -w -l .

simplify:
	gofmt -w -l -s .

check:
	go test -v -test.timeout=1200s $(PROJECT)/...

install:
	go install -v $(PROJECT)/...

clean:
	rm -rf bin/ || true
	go clean $(PROJECT)/...

all: clean check build_all

build_all: bin/tmp/$(EXE) bin/tmp/$(EXE) $(COMPRESSED_EXE_TARGETS)

release: clean bin/tmp/$(EXE) bin/tmp/$(EXE) $(COMPRESSED_EXE_TARGETS)
	github-release release \
    	--user ${USER} \
    	--repo ${REPO} \
    	--tag v${RELEASE_VERSION} \
    	--name "${NAME} v${RELEASE_VERSION}" \
    	--description "$(RELEASE_DESC_READ_CMD)" \
    	--pre-release
	$(foreach FILE,$(COMPRESSED_EXE),$(UPLOAD_CMD);)

docker: clean test build_all
	docker build --no-cache -t registry.kedocloud-internal.nl/silo:${VERSION} .

.PHONY: build check install clean format simplify build_all release

