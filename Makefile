OUT := wrench-turn 
PKG := github.com/okdv/wrench-turn 
# get version from git data
VERSION := v1.2.3-alpha

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
# commit and tag changes
	git checkout -b release-${VERSION}
	git add frontend/package.json version/version.go 
	git commit -m "Bump version to ${VERSION} for release" 
	git tag -a ${VERSION} -m "${VERSION}"
	git push origin release-${VERSION}
	git push --tags