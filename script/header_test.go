package script

import (
	"reflect"
	"regexp"
	"testing"
)

func TestDateRegex(t *testing.T) {
	re := regexp.MustCompile("^" + dateRegex + "$")
	if !re.Match([]byte(startFormat)) {
		t.Fatal("Regex doesn't match the date format")
	}
}

func TestParse(t *testing.T) {
	testCases := []struct {
		str          string
		expectedZone int
		expected     Header
	}{
		{
			str:          `Script started on 2025-04-01 12:34:56-07:00 [TERM="xterm-256color" TTY="/dev/pts/2" COLUMNS="155" LINES="25"]`,
			expectedZone: -7 * 60 * 60,
			expected: Header{
				Term:    "xterm-256color",
				Tty:     "/dev/pts/2",
				Columns: 155,
				Lines:   25,
			},
		},
		{ // parsing of ambiguous header should maximize command
			str:          `Script started on 2025-04-01 12:34:56+05:00 [COMMAND="echo \" TERM=""test" TERM="xterm-256color" TTY="/dev/pts/1" COLUMNS="214" LINES="25"]`,
			expectedZone: 5 * 60 * 60,
			expected: Header{
				Command: `echo \" TERM=""test`,
				Term:    "xterm-256color",
				Tty:     "/dev/pts/1",
				Columns: 214,
				Lines:   25,
			},
		},
	}

	for i := range testCases {
		h, err := ParseHeader(testCases[i].str)
		if err != nil {
			t.Fatalf("Test %d: Unexpected error: %v", i, err)
		}

		testCases[i].expected.Start = h.Start // todo: test that time zone has correct offset
		_, zone := h.Start.Zone()
		if zone != testCases[i].expectedZone { // todo: Make test parameter
			t.Fatalf("Wrong timezone (expected 7, got %d)", zone)
		}

		if !reflect.DeepEqual(h, testCases[i].expected) {
			t.Fatalf("Test %d:\nExpected: %#v\nActual: %#v", i, testCases[i].expected, h)
		}
	}
}
