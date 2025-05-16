// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package script

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

func TestRoundtrip(t *testing.T) {
	events := []Event{
		Event{
			Data:           "hello",
			ElapsedSeconds: 1.23,
			Code:           'O',
		},
		Event{
			Data:           "word",
			ElapsedSeconds: 1.23,
			Code:           'I',
		},
	}

	var typescript, timing bytes.Buffer

	for i, event := range events {
		if err := event.WriteAdvanced(&typescript, &timing); err != nil {
			t.Fatalf("Error writing event %d: %v", i, err)
		}
	}

	timingBuffered := bufio.NewReader(&timing)
	eventsOut := make([]Event, len(events))
	for i := range eventsOut {
		if err := eventsOut[i].Take(&typescript, timingBuffered); err != nil {
			t.Fatalf("Error reading event %d: %v", i, err)
		}
	}

	if !reflect.DeepEqual(events, eventsOut) {
		t.Errorf("Decoded events not the same as original:\nExpected: %#v\nActual:   %#v\n", events, eventsOut)
	}
}
