package grpchelper

import (
	"context"
	"errors"

	"github.com/andrescosta/goico/pkg/broadcaster"
	pb "github.com/andrescosta/jobico/api/types"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	ErrListeningData  = errors.New("error")
	ErrPublishingData = errors.New("error")
)

type GrpcBroadcaster[T, S proto.Message] struct {
	b *broadcaster.Broadcaster[T]
}

func StartBroadcaster[T, S proto.Message](ctx context.Context) *GrpcBroadcaster[T, S] {
	return &GrpcBroadcaster[T, S]{
		b: broadcaster.Start[T](ctx),
	}
}

func (b *GrpcBroadcaster[T, S]) Stop() {
	b.b.Stop()
}

func (b *GrpcBroadcaster[T, S]) Broadcast(_ context.Context, value S, utype pb.UpdateType) {
	var prototype T
	n := b.new(prototype, value, utype)
	b.b.Write(n)
}

func (b *GrpcBroadcaster[T, S]) RcvAndDispatchUpdates(ctx context.Context, s grpc.ServerStream) error {
	l, err := b.b.Subscribe()
	if err != nil {
		return err
	}
	logger := zerolog.Ctx(ctx)
	for {
		select {
		case <-ctx.Done():
			b.b.Unsubscribe(l)
			return ctx.Err()
		case <-s.Context().Done():
			b.b.Unsubscribe(l)
			return s.Context().Err()
		case d, ok := <-l.C:
			if !ok {
				return ErrListeningData
			}
			err := s.SendMsg(d)
			if err != nil {
				logger.Err(err).Msg("error sending data")
				return ErrPublishingData
			}
		}
	}
}

func (b *GrpcBroadcaster[T, S]) new(prototype T, value S, utype pb.UpdateType) T {
	v := prototype.ProtoReflect().New()
	o := v.Descriptor().Fields().ByName("object")
	t := v.Descriptor().Fields().ByName("type")
	v.Set(o, protoreflect.ValueOf(value.ProtoReflect()))
	v.Set(t, protoreflect.ValueOfEnum(utype.Number()))
	res, _ := v.Interface().(T)
	return res
}
