# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM golang:1.26.8-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /data

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN echo "Building for ${TARGETOS}/${TARGETARCH}" \
    && mkdir -p bin \
    && CGO_ENABLED=0 \
       GOOS="${TARGETOS}" \
       GOARCH="${TARGETARCH}" \
       go build -tags timetzdata -o bin/cva ./cmd/main.go

FROM scratch

WORKDIR /data

COPY --from=builder /data/bin/cva /data/cva
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /data/docs/openAPI /data/docs/docs/openAPI

ENV TZ="Europe/Berlin"

CMD ["/data/cva"]