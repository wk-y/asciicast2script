// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package script

import (
	"bufio"
	"fmt"
	"io"
)

type Event struct {
	Data           string
	ElapsedSeconds float64
	Code           rune
}

func (e *Event) Take(typescript io.Reader, timingfile *bufio.Reader) error {
	line, err := timingfile.ReadBytes('\n')
	if err != nil {
		return err
	}

	if len(line) == 0 {
		return fmt.Errorf("event line empty")
	}

	var dataLen int
	if line[0] >= 'A' && line[0] <= 'Z' {
		e.Code, e.ElapsedSeconds, dataLen, err = parseAdvancedTiming(string(line))
	} else {
		e.Code = 'O'
		e.ElapsedSeconds, dataLen, err = parseClassicTiming(string(line))
	}
	if err != nil {
		return err
	}

	if dataLen < 0 {
		return fmt.Errorf("negative event length in timing file")
	}

	buf := make([]byte, dataLen)
	for read := 0; read < dataLen; {
		n, err := typescript.Read(buf[read:])
		read += n
		if err != nil {
			return err
		}
	}

	e.Data = string(buf)

	return nil
}

func parseAdvancedTiming(s string) (code rune, elapsed float64, dataLen int, err error) {
	// todo: error if string continues past fields
	_, err = fmt.Sscanf(s, "%c %f %d", &code, &elapsed, &dataLen)
	return
}

func parseClassicTiming(s string) (elapsed float64, dataLen int, err error) {
	// todo: error if string continues past fields
	_, err = fmt.Sscanf(s, "%f %d", &elapsed, &dataLen)
	return
}

// Write event in advanced timing format
func (e *Event) WriteAdvanced(typescript, timingfile io.Writer) error {
	if _, err := fmt.Fprintf(timingfile, "%c %f %d\n", e.Code, e.ElapsedSeconds, len(e.Data)); err != nil {
		return err
	}

	if _, err := typescript.Write([]byte(e.Data)); err != nil {
		return err
	}

	return nil
}
