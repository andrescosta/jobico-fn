package test

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/andrescosta/goico/pkg/service"
	rpc "google.golang.org/grpc"
)

type GrpcConn struct {
	dialer   service.GrpcDialer
	listener service.GrpcListener
}
type HTTPConn struct {
	httpListener    service.HTTPListener
	httpTransporter service.HTTPTranporter
	clientBuilder   service.HTTPClientBuilder
}

func (g GrpcConn) Dial(ctx context.Context, addr string) (*rpc.ClientConn, error) {
	return g.dialer.Dial(ctx, addr)
}

func (g GrpcConn) Listen(addr string) (net.Listener, error) {
	return g.listener.Listen(addr)
}

func (h HTTPConn) Listen(addr string) (net.Listener, error) {
	return h.httpListener.Listen(addr)
}

func (h HTTPConn) Tranport(addr string) (*http.Transport, error) {
	return h.httpTransporter.Tranport(addr)
}

func (h HTTPConn) NewHTTPClient(addr string) (*http.Client, error) {
	return h.clientBuilder.NewHTTPClient(addr)
}

type Conn struct {
	conn net.Conn
}

func (c Conn) Read(b []byte) (n int, err error) {
	return c.conn.Read(b)
}

func (c Conn) Write(b []byte) (n int, err error) {
	return c.conn.Write(b)
}

func (c Conn) Close() error {
	return c.conn.Close()
}

func (c Conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c Conn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c Conn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c Conn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
