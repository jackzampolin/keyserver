FROM golang:alpine AS build-env

# Injest build args from Makefile
ARG BINARY=go-boilerplate
ARG GITHUB_USERNAME=myusername
ARG GOARCH=amd64

# Set up dependencies
ENV PACKAGES make git curl

# Set working directory for the build
WORKDIR /go/src/github.com/${GITHUB_USERNAME}/${BINARY}

# Install dependencies
RUN apk add --update $PACKAGES

# Install dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# Add source files
COPY . .

# Make the binary
RUN make install

# Final image
FROM alpine:edge

ARG BINARY=go-boilerplate
ARG GITHUB_USERNAME=myusername
ARG GOARCH=amd64
ENV BINARY=${BINARY}

# Install ca-certificates
RUN apk add --update ca-certificates
WORKDIR /root

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/${BINARY} /usr/bin/${BINARY}

# Run ${BINARY} by default
CMD ${BINARY}
