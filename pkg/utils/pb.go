package utils

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Rewire a protocol buffer message into a generic map[string]interface{} as an
// intermediate step before JSON or YAML marshalling. This is typically unnecessary work
// and is used as a workaround for multi-protocol systems.
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
