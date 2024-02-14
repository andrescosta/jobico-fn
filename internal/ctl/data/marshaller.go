package data

import (
	"google.golang.org/protobuf/proto"
)

type ProtoMessageMarshaller struct {
	prototype proto.Message
}

func (s *ProtoMessageMarshaller) Unmarshal(d []byte) (proto.Message, error) {
	i := s.prototype.ProtoReflect().New().Interface()
	if err := proto.Unmarshal(d, i); err != nil {
		return nil, err
	}
	return i, nil
}

func (s *ProtoMessageMarshaller) Marshal(q proto.Message) (string, []byte, error) {
	f := q.ProtoReflect().Descriptor().Fields().ByName("ID")
	id := q.ProtoReflect().Get(f).String()
	r, err := proto.Marshal(q)
	return id, r, err
}
