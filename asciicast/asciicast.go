// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package asciicast

import (
	"encoding/json"
	"fmt"
)

type Event struct {
	Time float64
	Code string
	Data string
}

var _ json.Marshaler = Event{}
var _ json.Unmarshaler = &Event{}

func (e Event) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{e.Time, e.Code, e.Data})
}

func (e *Event) UnmarshalJSON(data []byte) error {
	var err error
	var msg []any

	err = json.Unmarshal(data, &msg)
	if err != nil {
		return err
	}

	if len(msg) != 3 {
		return fmt.Errorf("expected 3 fields in event, got %d", len(msg))
	}

	var ok bool
	e.Time, ok = msg[0].(float64)
	if !ok {
		return fmt.Errorf("wrong type for event time field")
	}

	e.Code, ok = msg[1].(string)
	if !ok {
		return fmt.Errorf("wrong type for event code field")
	}

	e.Data, ok = msg[2].(string)
	if !ok {
		return fmt.Errorf("wrong type for event data field")
	}

	return nil
}
