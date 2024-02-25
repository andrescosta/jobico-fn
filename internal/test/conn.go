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
	listener      service.HTTPListener
	transporter   service.HTTPTranporter
	clientBuilder service.HTTPClientBuilder
}

func (g GrpcConn) Dial(ctx context.Context, addr string) (*rpc.ClientConn, error) {
	return g.dialer.Dial(ctx, addr)
}

func (g GrpcConn) Listen(addr string) (net.Listener, error) {
	listener, err := g.listener.Listen(addr)
	if err != nil {
		return nil, err
	}
	return &HTTPListener{
		listener: listener,
	}, nil
}

func (h HTTPConn) Listen(addr string) (net.Listener, error) {
	listener, err := h.listener.Listen(addr)
	if err != nil {
		return nil, err
	}
	return &HTTPListener{
		listener: listener,
	}, nil
}

func (h HTTPConn) Tranport(addr string) (*http.Transport, error) {
	return h.transporter.Tranport(addr)
}

func (h HTTPConn) NewHTTPClient(addr string) (*http.Client, error) {
	return h.clientBuilder.NewHTTPClient(addr)
}

type HTTPListener struct {
	listener net.Listener
}

func (l *HTTPListener) Accept() (net.Conn, error) { return l.listener.Accept() }

func (l *HTTPListener) Close() error { return l.listener.Close() }

func (l *HTTPListener) Addr() net.Addr { return l.listener.Addr() }

type GrpcListener struct {
	listener net.Listener
}

func (l *GrpcListener) Accept() (net.Conn, error) { return l.listener.Accept() }

func (l *GrpcListener) Close() error { return l.listener.Close() }

func (l *GrpcListener) Addr() net.Addr { return l.listener.Addr() }

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
