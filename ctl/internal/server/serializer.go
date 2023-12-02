package server

import (
	"github.com/andrescosta/goico/pkg/reflectico"
	"google.golang.org/protobuf/proto"
)

type newProtoMessage func() proto.Message

type ProtoMessageMarshaler struct {
	newMessage newProtoMessage
}

func (s *ProtoMessageMarshaler) Marshal(q proto.Message) ([]byte, error) {
	return proto.Marshal(q)
}

func (s *ProtoMessageMarshaler) Unmarshal(d []byte) (proto.Message, error) {
	dd := s.newMessage()
	if err := proto.Unmarshal(d, dd); err != nil {
		return nil, err
	}
	return dd, nil
}

func (s *ProtoMessageMarshaler) MarshalObj(q proto.Message) (string, []byte, error) {
	id := reflectico.GetFieldString(q, "ID")
	r, err := proto.Marshal(q)
	return id, r, err

}
