package types

import "encoding/json"

type Event struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

func NewEvent(event_type string, payload interface{}) (*Event, error) {

	event := &Event{
		Type: event_type,
	}

	if payload != nil {

		data, err := json.Marshal(payload)

		if err != nil {
			return nil, err
		}

		event.Data = data

	}

	return event, nil

}

func (event *Event) Unmarshal(target interface{}) error {

	if len(event.Data) == 0 {
		return nil
	}

	return json.Unmarshal(event.Data, target)

}
