FROM golang:1.11.1-alpine AS builder
RUN apk update && apk add --no-cache git gcc musl-dev
WORKDIR /build
COPY . .
RUN go install -v

FROM alpine
RUN apk update && apk add --no-cache ca-certificates
EXPOSE 3000
COPY --from=builder /go/bin/slack-docker /
CMD ["/slack-docker"]
