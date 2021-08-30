package rerpc

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Codec defines the interface reRPC uses to encode and decode messages.
type Codec interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

type jsonProtobufCodec struct {
	marshaler   protojson.MarshalOptions
	unmarshaler protojson.UnmarshalOptions
}

func (c jsonProtobufCodec) Marshal(value interface{}) ([]byte, error) {
	protoMessage, ok := value.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("could not case %T to a proto.Message", value)
	}
	return c.marshaler.Marshal(protoMessage)
}

func (c jsonProtobufCodec) Unmarshal(data []byte, value interface{}) error {
	protoMessage, ok := value.(proto.Message)
	if !ok {
		return fmt.Errorf("could not case %T to a proto.Message", value)
	}
	return c.unmarshaler.Unmarshal(data, protoMessage)
}

type protobufCodec struct{}

func (protobufCodec) Marshal(value interface{}) ([]byte, error) {
	protoMessage, ok := value.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("could not case %T to a proto.Message", value)
	}
	return proto.Marshal(protoMessage)
}

func (protobufCodec) Unmarshal(data []byte, value interface{}) error {
	protoMessage, ok := value.(proto.Message)
	if !ok {
		return fmt.Errorf("could not case %T to a proto.Message", value)
	}
	return proto.Unmarshal(data, protoMessage)
}
