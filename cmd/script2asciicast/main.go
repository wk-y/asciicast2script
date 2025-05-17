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

	"github.com/wk-y/asciicast2script/asciicast"
	"github.com/wk-y/asciicast2script/script"
)

var typescriptPath string
var timingfilePath string
var overwrite bool
var v3 bool // write asciicast v3

func init() {
	flag.StringVar(&typescriptPath, "typescript", "typescript", "input typescript file")
	flag.StringVar(&timingfilePath, "timingfile", "timingfile", "input timing file")
	flag.BoolVar(&overwrite, "overwrite", false, "overwrite existing output file")
	flag.BoolVar(&v3, "v3", false, "use asciicast v3 format")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTION]... OUTFILE.cast\n\n", os.Args[0])
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

	outFlags := os.O_WRONLY | os.O_CREATE
	if !overwrite {
		outFlags |= os.O_EXCL
	} else {
		outFlags |= os.O_TRUNC
	}

	cast := os.Stdin
	if castFile != "-" {
		var err error
		cast, err = os.OpenFile(castFile, outFlags, 0644)
		if err != nil {
			panic(err)
		}
		defer cast.Close()
	}

	script, err := os.Open(typescriptPath)
	if err != nil {
		panic(err)
	}
	defer script.Close()

	timing, err := os.Open(timingfilePath)
	if err != nil {
		panic(err)
	}
	defer timing.Close()

	converter := scriptToAsciicast
	if v3 {
		converter = scriptToAsciicastV3
	}

	err = converter(script, bufio.NewReader(timing), cast)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func scriptToAsciicast(typescript io.Reader, timingfile *bufio.Reader, cast io.Writer) error {
	encoder := json.NewEncoder(cast)
	tsBuffered := bufio.NewReader(typescript)

	headerBytes, err := tsBuffered.ReadBytes('\n')
	if err != nil {
		return err
	}

	header, err := script.ParseHeader(string(headerBytes))
	if err != nil {
		return err
	}

	timestamp := header.Start.Unix()
	env := map[string]string{}
	if header.Term != "" {
		env["TERM"] = header.Term
	}

	acHeader := asciicast.HeaderV2{
		Version:   2,
		Width:     header.Columns,
		Height:    header.Lines,
		Timestamp: &timestamp,
		Env:       env,
	}

	if header.Command != "" {
		acHeader.Command = &header.Command
	}

	if err := encoder.Encode(acHeader); err != nil {
		return err
	}

	var sEvent script.Event
	var time float64
	for {
		if err := sEvent.Take(tsBuffered, timingfile); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		time += sEvent.ElapsedSeconds

		acEvent := asciicast.Event{
			Time: time,
			Data: sEvent.Data,
		}

		var ignore bool
		switch sEvent.Code {
		case 'I':
			acEvent.Code = "i"
		case 'O':
			acEvent.Code = "o"
		default:
			ignore = true
		}

		if ignore {
			continue
		}

		if err := encoder.Encode(acEvent); err != nil {
			return err
		}
	}
}

func scriptToAsciicastV3(typescript io.Reader, timingfile *bufio.Reader, cast io.Writer) error {
	encoder := json.NewEncoder(cast)
	tsBuffered := bufio.NewReader(typescript)

	headerBytes, err := tsBuffered.ReadBytes('\n')
	if err != nil {
		return err
	}

	header, err := script.ParseHeader(string(headerBytes))
	if err != nil {
		return err
	}

	timestamp := header.Start.Unix()
	env := map[string]string{}
	if header.Term != "" {
		env["TERM"] = header.Term
	}

	acHeader := asciicast.HeaderV3{
		Version: 3,
		Term: asciicast.TermInfo{
			Cols: header.Columns,
			Rows: header.Lines,
			Type: &header.Term,
		},
		Timestamp: &timestamp,
		Env:       env,
	}

	if header.Command != "" {
		acHeader.Command = &header.Command
	}

	if err := encoder.Encode(acHeader); err != nil {
		return err
	}

	var sEvent script.Event
	for {
		if err := sEvent.Take(tsBuffered, timingfile); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		acEvent := asciicast.Event{
			Time: sEvent.ElapsedSeconds, // asciicast v3 uses relative time
			Data: sEvent.Data,
		}

		var ignore bool
		switch sEvent.Code {
		case 'I':
			acEvent.Code = "i"
		case 'O':
			acEvent.Code = "o"
		default:
			ignore = true
		}

		if ignore {
			continue
		}

		if err := encoder.Encode(acEvent); err != nil {
			return err
		}
	}
}
