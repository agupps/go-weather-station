ARG DEPENDENCY_PROXY
FROM --platform=$BUILDPLATFORM ${DEPENDENCY_PROXY}golang:1.21 as builder

WORKDIR /workspace

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# Cache deps before building and copying source so that we don't need to re-download as much and so that source changes
# don't invalidate our downloaded layer
RUN go mod download

# Copy the Go source
COPY cmd/main.go cmd/main.go
COPY internal/ internal/

# Build
# the GOARCH has not a default value to allow the binary be built according to the host where the command
# was called. For example, if we call make docker-build in a local env which has the Apple Silicon M1 SO
# the docker BUILDPLATFORM arg will be linux/arm64 when for Apple x86 it will be linux/amd64. Therefore,
# by leaving it empty we can ensure that the container and binary shipped on it will have the same platform.
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -a -o /go-weather-station cmd/main.go

# Use distroless as minimal base image to package the handler binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /go-weather-station .
USER 65532:65532

# checkov:skip=CKV_DOCKER_2: Adding a HEALTHCHECK instruction is unnecessary given that the image is only run in Kubernetes
# HEALTHCHECK

ENTRYPOINT ["/go-weather-station"]


