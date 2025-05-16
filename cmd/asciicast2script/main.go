// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/wk-y/asciicast2script/asciicast"
	"github.com/wk-y/asciicast2script/script"
)

var typescriptPath string
var timingfilePath string
var overwrite bool

func init() {
	flag.StringVar(&typescriptPath, "typescript", "typescript", "output typescript file")
	flag.StringVar(&timingfilePath, "timingfile", "timingfile", "output timing file")
	flag.BoolVar(&overwrite, "overwrite", false, "overwrite existing files")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTION]... ASCIICAST\n\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	argv := flag.Args()
	if len(argv) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	castFile := argv[0]

	cast := os.Stdin
	if castFile != "-" {
		var err error
		cast, err = os.Open(castFile)
		if err != nil {
			panic(err)
		}
		defer cast.Close()
	}

	outFlags := os.O_WRONLY | os.O_CREATE
	if !overwrite {
		outFlags |= os.O_EXCL
	}

	script, err := os.OpenFile(typescriptPath, outFlags, 0644)
	if err != nil {
		panic(err)
	}
	defer script.Close()

	timing, err := os.OpenFile(timingfilePath, outFlags, 0644)
	if err != nil {
		panic(err)
	}
	defer timing.Close()

	err = asciicastToScript(cast, script, timing)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func asciicastToScript(cast io.Reader, typescript, timingfile io.Writer) error {
	buffered := bufio.NewReader(cast)

	// Convert the header line
	headerBytes, err := buffered.ReadBytes('\n')
	if err != nil {
		return err
	}

	header, err := asciicast.DecodeHeader(headerBytes)
	if err != nil {
		return err
	}

	fmt.Fprintln(typescript, asciicastHeaderToScript(header))

	// Convert events
	var previousEventTime float64
	for {
		line, err := buffered.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		if strings.HasPrefix(string(line), "#") { // comment line
			// todo: error if comment encountered in v2 file?
			continue
		}

		var acEvent asciicast.Event
		if err := json.Unmarshal(line, &acEvent); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		sEvent := script.Event{
			Data:           acEvent.Data,
			ElapsedSeconds: acEvent.Time - previousEventTime,
		}

		var ignore bool
		switch acEvent.Code {
		case "i":
			sEvent.Code = 'I'
		case "o":
			sEvent.Code = 'O'
		default:
			ignore = true
		}

		if ignore {
			continue
		}

		if err := sEvent.WriteAdvanced(typescript, timingfile); err != nil {
			return err
		}

		previousEventTime = acEvent.Time
	}
}

func asciicastHeaderToScript(header asciicast.Header) script.Header {
	var result script.Header
	if timestamp, ok := header.Timestamp(); ok {
		result.Start = time.Unix(timestamp, 0)
	}

	if term, ok := header.Term(); ok {
		result.Term = term
	}

	result.Columns = header.Width()
	result.Lines = header.Height()
	return result
}
