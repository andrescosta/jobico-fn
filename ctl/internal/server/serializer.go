package server

import (
	"google.golang.org/protobuf/proto"
)

type ProtoMessageMarshaler struct {
	prototype proto.Message
}

func (s *ProtoMessageMarshaler) Marshal(q proto.Message) ([]byte, error) {
	return proto.Marshal(q)
}

func (s *ProtoMessageMarshaler) Unmarshal(d []byte) (proto.Message, error) {
	i := s.prototype.ProtoReflect().New().Interface()
	if err := proto.Unmarshal(d, i); err != nil {
		return nil, err
	}
	return i, nil
}

func (s *ProtoMessageMarshaler) MarshalObj(q proto.Message) (string, []byte, error) {
	f := q.ProtoReflect().Descriptor().Fields().ByName("ID")
	id := q.ProtoReflect().Get(f).String()
	r, err := proto.Marshal(q)
	return id, r, err

}
