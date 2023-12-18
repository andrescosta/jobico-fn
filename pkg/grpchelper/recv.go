package grpchelper

import (
	"github.com/andrescosta/goico/pkg/broadcaster"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func Recv[T proto.Message](s grpc.ClientStream, c chan<- T) error {
	logger := zerolog.Ctx(s.Context())

	for {
		select {
		case <-s.Context().Done():
			if err := s.CloseSend(); err != nil {
				logger.Warn().AnErr("error", err).Msg("Recv: Error while closing stream.")
			}
			return s.Context().Err()
		default:
			var t T
			p := t.ProtoReflect().New()
			err := s.RecvMsg(p.Interface())
			if err != nil {
				if status.Code(err) != codes.Canceled {
					logger.Warn().AnErr("error", err).Msg("Recv: error getting message")
				}
				continue
			}
			select {
			case c <- p.Interface().(T):
			default:
			}
		}
	}
}

func Listen[T proto.Message](s grpc.ClientStream, b *broadcaster.Broadcaster[T]) error {
	logger := zerolog.Ctx(s.Context())

	for {
		select {
		case <-s.Context().Done():
			if err := s.CloseSend(); err != nil {
				logger.Warn().AnErr("error", err).Msg("Listen: Error while closing stream.")
			}
			return s.Context().Err()
		default:
			var t T
			p := t.ProtoReflect().New()
			err := s.RecvMsg(p.Interface())
			if err != nil {
				if status.Code(err) != codes.Canceled {
					logger.Warn().AnErr("error", err).Msg("Listen: error getting message")
				}
				continue
			}
			b.Write(p.Interface().(T))
		}
	}
}
