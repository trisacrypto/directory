package wire

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Rewire a protocol buffer message into a generic map[string]interface{} as an
// intermediate step before JSON or YAML marshalling. This is typically unnecessary work
// and is used as a workaround for multi-protocol systems.
// This method is primarily being used by the Admin API to convert protocol buffer
// messages into JSON generics for serialization by Gin.
func Rewire(m protoreflect.ProtoMessage) (out map[string]interface{}, err error) {
	// Serialize the VASP from protojson
	jsonpb := protojson.MarshalOptions{
		Multiline:       false,
		AllowPartial:    true,
		UseProtoNames:   true,
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
	}

	var data []byte
	if data, err = jsonpb.Marshal(m); err != nil {
		return nil, err
	}

	// Remarshal the JSON (unnecessary work, but done to make things easier)
	if err = json.Unmarshal(data, &out); err != nil {
		return nil, err
	}

	return out, nil
}

// Unwire a generic map[string]interface{} into a protocol buffer message as an
// intermediate step to parsing a JSON or YAML request instead of protocol buffers. Like
// Rewire, this is a workaround for multi-protocol systems.
// This method is primarily being used by the Admin API to convert JSON generics
// unmarshaled by Gin into protocol buffers for interacting with the data store.
func Unwire(entry map[string]interface{}, msg protoreflect.ProtoMessage) (err error) {
	// Serialize the data into JSON format
	var data []byte
	if data, err = json.Marshal(entry); err != nil {
		return err
	}

	// Remarshal the JSON into a protocol buffer via protojson
	jsonpb := protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	if err = jsonpb.Unmarshal(data, msg); err != nil {
		return err
	}
	return nil
}
