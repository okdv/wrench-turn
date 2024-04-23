OUT := wrench-turn 
PKG := github.com/okdv/wrench-turn 
# get version from git data
VERSION := $(shell git describe --tags --abbrev=0)

.PHONY: build 

build:
# set version in package.json
	sed -i 's/"version": "[^"]*"/"version": "${VERSION}"/' frontend/package.json
# update version in go 
	sed -i "s/Version = \".*\"/Version = \"${VERSION}\"/g" version/version.go 
# run go test 
	@go test 
# build go app
	go build -v -o ${OUT} -ldflags="-X main.version=${VERSION}" ${PKG} 