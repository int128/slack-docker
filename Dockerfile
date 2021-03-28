FROM golang:1.16 as builder

WORKDIR /builder
COPY go.* .
RUN go mod download
COPY Makefile .
COPY main.go .
COPY cmd cmd
COPY formatter formatter
RUN make

FROM gcr.io/distroless/base-debian10
COPY --from=builder /builder/slack-docker /
CMD ["/slack-docker"]
