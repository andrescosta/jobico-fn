package server

import (
	"strconv"

	"github.com/andrescosta/goico/pkg/reflectico"
	"google.golang.org/protobuf/proto"
)

type newProtoMessage func() proto.Message

type ProtoMessageMarshaler struct {
	newMessage newProtoMessage
}

func (s *ProtoMessageMarshaler) Marshal(id uint64, q proto.Message) ([]byte, error) {
	reflectico.SetFieldString(q, "ID", strconv.FormatUint(id, 10))
	return proto.Marshal(q)
}

func (s *ProtoMessageMarshaler) Unmarshal(id uint64, d []byte) (proto.Message, error) {
	dd := s.newMessage()
	if err := proto.Unmarshal(d, dd); err != nil {
		return nil, err
	}
	return dd, nil
}

func (s *ProtoMessageMarshaler) MarshalObj(q proto.Message) (uint64, []byte, error) {
	id := reflectico.GetFieldUInt(q, "ID")
	r, err := proto.Marshal(q)
	return id, r, err

}
