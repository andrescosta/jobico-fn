FROM golang:1.21.0-bullseye as builder

ENV CGO_CPPFLAGS="-D_FORTIFY_SOURCE=2 -fstack-protector-all"
ENV GOFLAGS="-buildmode=pie"

WORKDIR /workdir

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -ldflags "-s -w" -trimpath ./cmd/...

# try USER with nobody:nobody 
FROM gcr.io/distroless/base-debian11:nonroot as ctl
COPY --from=builder /workdir/ctl /bin/ctl
USER 65534:65534
CMD ["/bin/ctl"]

FROM gcr.io/distroless/base-debian11:nonroot as exec
COPY --from=builder /workdir/executor /bin/executor
USER 65534:65534
CMD ["/bin/executor"]

FROM gcr.io/distroless/base-debian11:nonroot as listener
COPY --from=builder /workdir/listener /bin/listener
USER 65534:65534
CMD ["/bin/listener"]

FROM gcr.io/distroless/base-debian11:nonroot as queue
COPY --from=builder /workdir/queue /bin/queue
USER 65534:65534
CMD ["/bin/queue"]

FROM gcr.io/distroless/base-debian11:nonroot as recorder
COPY --from=builder /workdir/recorder /bin/recorder
USER 65534:65534
CMD ["/bin/recorder"]

FROM gcr.io/distroless/base-debian11:nonroot as repo
COPY --from=builder /workdir/repo /bin/repo
USER 65534:65534
CMD ["/bin/repo"]