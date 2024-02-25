package test

import (
	"reflect"
	"time"
)

type Error int

const (
	ErrTimeout Error = iota
	ErrSlowOp
	ErrConnect
)

type (
	serviceWithError struct {
		service any
		errors  []any
	}
	PlatformBuilder struct {
		curr     any
		services map[reflect.Type]serviceWithError
	}

	errorConn struct {
		retries int
	}
	errorReadtimeout struct{}

	errorPlat struct {
		typ error
	}
)

func NewPlatformBuilder() *PlatformBuilder {
	s := make(map[reflect.Type]serviceWithError)

	return &PlatformBuilder{
		services: s,
	}
}

// Retries

func (p *PlatformBuilder) WithService(svc any) *PlatformBuilder {
	p.curr = svc
	return nil
}

func (p *PlatformBuilder) WithErrConn(retries int) *PlatformBuilder {
	return nil
}

func (p *PlatformBuilder) WithErrReadtimeout(d time.Duration) *PlatformBuilder {
	return nil
}

func (p *PlatformBuilder) WithErrWritetimeout(d time.Duration) *PlatformBuilder {
	return nil
}

func (p *PlatformBuilder) WithErrRead(n int) *PlatformBuilder {
	return nil
}

func (p *PlatformBuilder) WithErrWrite(n int) *PlatformBuilder {
	return nil
}

func (p *PlatformBuilder) WithErrClose(n int) *PlatformBuilder {
	return nil
}

func (p *PlatformBuilder) WithErrSlowClose(d time.Duration) *PlatformBuilder {
	return nil
}

func (p *PlatformBuilder) WithErrSlowRead(d time.Duration) *PlatformBuilder {
	return nil
}

func (p *PlatformBuilder) WithErrSlowWrite(d time.Duration) *PlatformBuilder {
	return nil
}

func (p *PlatformBuilder) WithErrSlowOp(d time.Duration) *PlatformBuilder {
	return nil
}
