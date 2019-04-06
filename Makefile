BINARY            = keyserver
GITHUB_USERNAME   = jackzampolin
DOCKER_REPO       = quay.io/jackzampolin
VERSION           = v0.1.0
GOARCH            = amd64
ARTIFACT_DIR      = build
PORT 							= 3000

COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Symlink into GOPATH
BUILD_DIR=${GOPATH}/src/github.com/${GITHUB_USERNAME}/${BINARY}
FLAG_PATH=github.com/${GITHUB_USERNAME}/${BINARY}/cmd
DOCKER_TAG=${VERSION}
DOCKER_IMAGE=${DOCKER_REPO}/${BINARY}

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X ${FLAG_PATH}.Version=${VERSION} -X ${FLAG_PATH}.Commit=${COMMIT} -X ${FLAG_PATH}.Branch=${BRANCH}"

# Build the project
all: clean linux darwin

# Build and Install project into GOPATH using current OS setup
install:
	go install ${LDFLAGS} ./...

test:
	go test -v ./api/...

# Build binary for Linux
linux: clean
	cd ${BUILD_DIR}; \
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${ARTIFACT_DIR}/${BINARY}-linux-${GOARCH} . ; \
	cd - >/dev/null

# Build binary for MacOS
darwin:
	cd ${BUILD_DIR}; \
	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ${ARTIFACT_DIR}/${BINARY}-darwin-${GOARCH} . ; \
	cd - >/dev/null

# Build binary for Windows
windows:
	cd ${BUILD_DIR}; \
	GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o ${ARTIFACT_DIR}/${BINARY}-darwin-${GOARCH} . ; \
	cd - >/dev/null

# Install golang dependencies

# Build the docker image and give it the appropriate tags
docker:
	cd ${BUILD_DIR} >/dev/null
	docker build \
		--build-arg BINARY=${BINARY} \
		--build-arg GITHUB_USERNAME=${GITHUB_USERNAME} \
		--build-arg GOARCH=${GOARCH} \
		-t ${DOCKER_IMAGE}:${DOCKER_TAG} \
		.
	docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:latest
	docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:${BRANCH}
	cd - >/dev/null

# Push the docker image to the configured repo
docker-push:
	cd ${BUILD_DIR} >/dev/null
	docker push ${DOCKER_IMAGE}:${DOCKER_TAG}
	docker push ${DOCKER_IMAGE}:${BRANCH}
	docker push ${DOCKER_IMAGE}:latest
	cd - >/dev/null

# Run the docker image as a server exposing the service port, mounting configuration from this repo
docker-run:
	cd ${BUILD_DIR} >/dev/null
	docker run -p ${PORT}:${PORT} -v ${BUILD_DIR}/${BINARY}.yaml:/root/.${BINARY}.yaml -it ${DOCKER_IMAGE}:${DOCKER_TAG} ${BINARY} serve
	cd - >/dev/null

# Remove all the built binaries
clean:
	cd ${BUILD_DIR} >/dev/null
	rm -rf ${ARTIFACT_DIR}/*
	cd - >/dev/null

.PHONY: linux darwin fmt clean
