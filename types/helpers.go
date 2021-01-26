package types

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func Unmarshal(data interface{}, raw []byte, legacy bool) (interface{}, error) {
	if !legacy {
		if msg, ok := data.(proto.Message); ok {
			uopts := protojson.UnmarshalOptions{
				DiscardUnknown: true,
			}

			if err := uopts.Unmarshal(raw, msg); err != nil {
				return nil, err
			}
			return msg, nil
		}
	}

	if err := json.Unmarshal(raw, data); err != nil {
		return nil, err
	}

	return data, nil
}

func Marshal(data interface{}, legacy bool) ([]byte, error) {
	if !legacy {
		if msg, ok := data.(proto.Message); ok {
			return protojson.Marshal(msg)
		}
	}

	return json.Marshal(data)
}
