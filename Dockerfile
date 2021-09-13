FROM golang:1.17 as builder

ARG VERSION

ENV CGO_ENABLED=0
WORKDIR /go/src/app
COPY . .
RUN go build -ldflags "-s -w -X github.com/kobtea/mixcli/cmd.Version=$VERSION"

FROM scratch
COPY --from=builder /go/src/app/mixcli /mixcli
ENTRYPOINT ["/mixcli"]
