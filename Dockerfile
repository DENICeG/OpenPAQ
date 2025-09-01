FROM golang:1.25.0-alpine3.21 as builder

RUN apk update && apk add --no-cache git ca-certificates

WORKDIR /data

RUN echo Building for linux
RUN mkdir -p bin
COPY . .
RUN go get -v -d ./...
RUN CGO_ENABLED=0 GOOS=linux go build -tags timetzdata -o bin/cva -a ./cmd/main.go

FROM scratch

WORKDIR /data
COPY --from=builder /data/bin/cva /data
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /data/docs/openAPI /data/docs/docs/openAPI
ENV TZ="Europe/Berlin"
CMD [ "/data/cva" ]

