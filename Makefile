BINARY = testapp
DOCKER_IMAGE = testapp
VET_REPORT = vet.report
TEST_REPORT = tests.xml
GOARCH = amd64

VERSION?=?
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.branch=${BRANCH}"

# Build the project
all: clean test vet linux darwin 

# Build linux binary
linux: 
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ./build/${BINARY}-linux-${GOARCH} main/main.go ; 

# Build darwin binary
darwin:
	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ./build/${BINARY}-darwin-${GOARCH} main/main.go ; 

# Run go tests on all packages
test: clean bootstrap
	go test ./...

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

.PHONY: linux darwin test vet fmt clean docker
