package events

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/contextgg/pkg/types"
	"github.com/pkg/errors"
)

// EventDecoder reader to event
func EventDecoder(reader io.Reader) (*Event, error) {
	out := struct {
		*Event

		Data json.RawMessage `json:"data"`
	}{}

	if err := json.NewDecoder(reader).Decode(&out); err != nil {
		return nil, errors.Wrapf(err, "Could not decode event")
	}

	typeData, ok := types.GetTypeData(out.Type)
	if !ok {
		return nil, fmt.Errorf("Could not find type with name %s", out.Type)
	}

	data := typeData.Factory()
	if err := json.Unmarshal(out.Data, data); err != nil {
		return nil, errors.Wrapf(err, "Could not decode event data")
	}

	evt := out.Event
	evt.Data = data

	if evt.Metadata == nil {
		evt.Metadata = make(map[string]interface{})
	}
	return evt, nil
}
