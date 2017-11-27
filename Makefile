BINARY = testapp
DOCKER_IMAGE = testapp
VET_REPORT = vet.report
TEST_REPORT = tests.xml
GOARCH = amd64
COVERAGEOUT := ./build/coverage.out
COVERAGETMP := ./build/coverage.tmp
PKGS := $(shell go list ./...)

VERSION?=?
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.branch=${BRANCH}"

# Build the project
all: clean bootstrap test vet linux darwin docker

# Build linux binary
linux: 
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ./build/${BINARY}-linux-${GOARCH} main/main.go ; 

# Build darwin binary
darwin:
	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ./build/${BINARY}-darwin-${GOARCH} main/main.go ; 

# Run go tests on all packages
test: clean
	@echo 'mode: atomic' > $(COVERAGEOUT) \
	exitcode=0; \
	for pkg in $(PKGS); do \
		go test -v -race -coverprofile=$(COVERAGETMP) -covermode=atomic $$pkg 2>&1 | grep -v 'warning: no packages being tested depend on'; \
		testexitcode=$$?; \
		if [ $$testexitcode -ne 0 ]; then \
			exitcode=$$testexitcode; \
		fi; \
		if [ -f $(COVERAGETMP) ]; then \
			grep -v -e 'mode: set' -e 'mode: atomic' $(COVERAGETMP) >> $(COVERAGEOUT); \
			rm $(COVERAGETMP); \
		fi; \
	done; \
	go tool cover -html=$(COVERAGEOUT) -o ./build/coverage.html; \
	exit $$exitcode

# Run go vet on all packages
vet:
	go vet ./... > ./build/${VET_REPORT} 2>&1 ;

# Clean the workspace
clean:
	rm -f ./build/*

# Download vendored deps
bootstrap:
	mkdir -p ./build
	dep ensure

# Build a docker image
docker: linux
	docker build -t ${DOCKER_IMAGE} .

prepare:
	mkdir -p ~/.docker_data/exampleapp

generate:
	mkdir -p ./mocks
	go generate ./...

.PHONY: linux darwin test vet fmt clean docker
