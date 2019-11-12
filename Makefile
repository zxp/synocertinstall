tag_name = -X main.tagName=$(shell git describe --abbrev=0 --tags)
branch = -X main.branch=$(shell git rev-parse --abbrev-ref HEAD)
commit_id = -X main.commitID=$(shell git log --pretty=format:"%h" -1)
build_time = -X main.buildTime=$(shell date -u '+%Y-%m-%d_%I:%M:%S%p')

VERSION = $(tag_name) $(branch) $(commit_id) $(build_time)

.PHONY: all
.DEFAULT_GOAL := osx_amd64

osx_amd64:
	env GOOS=darwin GOARCH=amd64 go build -v -ldflags "-s -w ${VERSION}" -o builds/synocertinstall_osx
	cd builds && 7z a synocertinstall_osx.7z synocertinstall_osx

linux: linux_amd64 

linux_amd64:
	env GOOS=linux GOARCH=amd64 go build -v -ldflags "-s -w ${VERSION}" -o builds/synocertinstall_lin
	cd builds && 7z a synocertinstall_lin.7z synocertinstall_lin

all: osx_amd64 linux 
