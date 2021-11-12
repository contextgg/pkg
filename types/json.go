package types

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type JsonOptions struct {
	Legacy bool
}
type JsonOption func(*JsonOptions)

func UseLegacyJsonSerializer(legacy bool) JsonOption {
	return func(o *JsonOptions) {
		o.Legacy = legacy
	}
}

func buildOpts(all []JsonOption) *JsonOptions {
	opts := &JsonOptions{}
	for _, a := range all {
		a(opts)
	}
	return opts
}

func EntryUnmarshal(entry *Entry, raw []byte, options ...JsonOption) (interface{}, error) {
	data := entry.Factory()
	if err := JsonUnmarshal(raw, data, options...); err != nil {
		return nil, err
	}
	return data, nil
}

func JsonUnmarshal(data []byte, v interface{}, options ...JsonOption) error {
	opts := buildOpts(options)

	if !opts.Legacy {
		if msg, ok := v.(proto.Message); ok {
			uopts := protojson.UnmarshalOptions{
				DiscardUnknown: true,
			}

			if err := uopts.Unmarshal(data, msg); err != nil {
				return err
			}
			return nil
		}
	}

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}

func JsonMarshal(data interface{}, options ...JsonOption) ([]byte, error) {
	opts := buildOpts(options)

	if !opts.Legacy {
		if msg, ok := data.(proto.Message); ok {
			return protojson.Marshal(msg)
		}
	}

	return json.Marshal(data)
}
