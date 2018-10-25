FROM golang:1.11.1-alpine AS builder
WORKDIR /go/src/github.com/int128/slack-docker
RUN apk update && apk add --no-cache git gcc musl-dev
ENV GO111MODULE on
COPY . .
RUN go install -v

FROM alpine
RUN apk update && apk add --no-cache ca-certificates
EXPOSE 3000
COPY --from=builder /go/bin/slack-docker /
CMD ["/slack-docker"]
