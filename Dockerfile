FROM golang:1.13.7-alpine AS builder
ENV CGO_ENABLED=0
WORKDIR /build
COPY . .
RUN go install -v

FROM alpine:3.11
EXPOSE 3000
COPY --from=builder /go/bin/slack-docker /
CMD ["/slack-docker"]
