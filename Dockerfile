FROM golang:1.21.5-bullseye as builder

ENV CGO_CPPFLAGS="-D_FORTIFY_SOURCE=2 -fstack-protector-all"
ENV GOFLAGS="-buildmode=pie"

WORKDIR /workdir

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg
COPY ./api ./api
RUN go build -o ./bin/ -ldflags "-s -w" -trimpath ./cmd/...

FROM  debian:12-slim as ctl
WORKDIR /app
COPY --from=builder /workdir/bin/ctl ctl
CMD ["/app/ctl","--env:basedir=/app"]

FROM  debian:12-slim as exec
COPY --from=builder /workdir/bin/executor /bin
ENTRYPOINT ["/bin/executor", "--env:basedir=/bin"]

FROM  debian:12-slim as listener
COPY --from=builder /workdir/bin/listener /bin/listener
ENTRYPOINT ["/bin/listener","--env:basedir=/bin"]

FROM  debian:12-slim as queue
COPY --from=builder /workdir/bin/queue /bin/queue
ENTRYPOINT ["/bin/queue","--env:basedir=/bin"]

FROM  debian:12-slim as recorder
COPY --from=builder /workdir/bin/recorder /bin/recorder
ENTRYPOINT ["/bin/recorder","--env:basedir=/bin"]

FROM  debian:12-slim as repo
COPY --from=builder /workdir/bin/repo /bin/repo
ENTRYPOINT ["/bin/repo","--env:basedir=/bin"]