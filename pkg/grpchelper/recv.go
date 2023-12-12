package grpchelper

import (
	"context"

	"github.com/andrescosta/goico/pkg/chico"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func Recv[T proto.Message](ctx context.Context, s grpc.ClientStream, c chan<- T) error {
	logger := zerolog.Ctx(ctx)
	for {
		select {
		case <-s.Context().Done():
			s.CloseSend()
			return nil
		case <-ctx.Done():
			s.CloseSend()
			return nil
		default:
			var t T
			p := t.ProtoReflect().New()
			err := s.RecvMsg(p.Interface())
			if err != nil {
				logger.Warn().Msgf("error getting message %s", err)
			} else {
				select {
				case c <- p.Interface().(T):
				default:
				}
			}
		}
	}
}

func Listen[T proto.Message](ctx context.Context, s grpc.ClientStream, b *chico.Broadcaster[T]) {
	logger := zerolog.Ctx(ctx)
	for {
		select {
		case <-s.Context().Done():
			s.CloseSend()
			return
		case <-ctx.Done():
			s.CloseSend()
			return
		default:
			var t T
			p := t.ProtoReflect().New()
			err := s.RecvMsg(p.Interface())
			if err != nil {
				logger.Warn().Msgf("error getting message %s", err)
			} else {
				b.Write(p.Interface().(T))
			}
		}
	}
}
