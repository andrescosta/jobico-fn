FROM golang:1.21.5-bullseye as builder

WORKDIR /workdir

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg
COPY ./api ./api
RUN go build -o ./bin/ -trimpath ./cmd/...

FROM  debian:12-slim as ctl
WORKDIR /app
COPY --from=builder /workdir/bin/ctl ctl
CMD ["/app/ctl","--env:basedir=/app"]

FROM  debian:12-slim as exec
WORKDIR /app
COPY --from=builder /workdir/bin/executor executor
ENTRYPOINT ["/app/executor", "--env:basedir=/app"]

FROM  debian:12-slim as listener
WORKDIR /app
COPY --from=builder /workdir/bin/listener listener
ENTRYPOINT ["/app/listener","--env:basedir=/app"]

FROM  debian:12-slim as queue
WORKDIR /app
COPY --from=builder /workdir/bin/queue queue
ENTRYPOINT ["/app/queue","--env:basedir=/app"]

FROM  debian:12-slim as recorder
WORKDIR /app
COPY --from=builder /workdir/bin/recorder recorder
ENTRYPOINT ["/app/recorder","--env:basedir=/app"]

FROM  debian:12-slim as repo
WORKDIR /app
COPY --from=builder /workdir/bin/repo repo
ENTRYPOINT ["/app/repo","--env:basedir=/app"]