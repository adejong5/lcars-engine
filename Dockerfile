# syntax=docker/dockerfile:1

# Build stage: cross-compile a static binary for the target arch. buildx sets
# TARGETOS/TARGETARCH; the stage itself runs on the native build platform for
# speed (Go cross-compiles without emulation).
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS build
RUN apk add --no-cache ca-certificates
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /out/lcars ./cmd/lcars

# Final stage: scratch + CA certs + the static binary (~10-15 MB, no shell).
# Same image runs as a standalone container or, with a thin manifest, as a
# multi-arch Home Assistant add-on (see decisions/0006).
FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /out/lcars /lcars
EXPOSE 8080
USER 65534:65534
ENTRYPOINT ["/lcars"]
