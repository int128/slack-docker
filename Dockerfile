FROM golang:1.18 as builder

WORKDIR /builder
COPY go.* .
RUN go mod download
COPY main.go .
COPY cmd cmd
COPY formatter formatter
RUN go build

FROM gcr.io/distroless/base-debian10
COPY --from=builder /builder/slack-docker /
CMD ["/slack-docker"]
