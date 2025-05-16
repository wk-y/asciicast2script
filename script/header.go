// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package script

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// time format of header start line
const startFormat = "2006-01-02 15:04:05-07:00"
const dateRegex = `[0-9]{4}\-[0-9]{2}\-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}[+-][0-9]{2}:[0-9]{2}`

type Header struct {
	Start   time.Time
	Command string
	Term    string
	Tty     string
	Columns int
	Lines   int
}

var _ fmt.Stringer = Header{}

func (h Header) String() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, `Script started on %s [`, h.Start.Format(startFormat))
	spacer := ""

	if h.Command != "" {
		fmt.Fprintf(&builder, `%sCOMMAND="%s"`, spacer, h.Command)
		spacer = " "
	}

	isTerminal := false

	if h.Term != "" {
		fmt.Fprintf(&builder, `%sTERM="%s"`, spacer, h.Term)
		spacer = " "
		isTerminal = true
	}

	if h.Tty != "" {
		fmt.Fprintf(&builder, `%sTTY="%s"`, spacer, h.Tty)
		spacer = " "
	}

	if h.Columns != 0 {
		fmt.Fprintf(&builder, `%sCOLUMNS="%d"`, spacer, h.Columns)
		spacer = " "
	}

	if h.Lines != 0 {
		fmt.Fprintf(&builder, `%sLINES="%d"`, spacer, h.Lines)
		spacer = " "
	}

	if !isTerminal {
		fmt.Fprintf(&builder, `%s<not executed on terminal>`, spacer)
	}

	fmt.Fprint(&builder, "]")

	return builder.String()
}

// Currently the regex assumes there will be no quotes in fields other than <command>.
// This is done to make quotes in the command match reliably, with the tradeoff that
// quotes in <term> or <tty> will cause those fields to be interpreted as part of the commmand.
var headerRegex = regexp.MustCompile("^Script started on (?P<date>" + dateRegex + `) \[` +
	`(?:COMMAND="(?P<command>.*?)" )?` +
	`(?:(?:` + // terminal exists
	`(?:TERM="(?P<term>[^"]*?)" )?` +
	`(?:TTY="(?P<tty>[^"]*?)" )?` +
	`COLUMNS="(?P<columns>[0-9]+)" LINES="(?P<lines>[0-9]+)"` +
	`)|(?:` + // not a terminal
	`<not executed on terminal>` +
	`))` +
	"")

func ParseHeader(header string) (result Header, err error) {
	match := headerRegex.FindSubmatch([]byte(header))
	if match == nil {
		return result, fmt.Errorf("improper header structure")
	}
	for i, subexp := range headerRegex.SubexpNames() {
		if match[i] == nil {
			continue
		}

		switch subexp {
		case "date":
			result.Start, err = time.ParseInLocation(startFormat, string(match[i]), time.UTC)
			if err != nil {
				return result, err
			}
		case "command":
			result.Command = string(match[i])
		case "term":
			result.Term = string(match[i])
		case "tty":
			result.Tty = string(match[i])
		case "columns":
			result.Columns, err = strconv.Atoi(string(match[i]))
			if err != nil {
				return result, err
			}
		case "lines":
			result.Lines, err = strconv.Atoi(string(match[i]))
			if err != nil {
				return result, err
			}
		}
	}
	return result, nil
}
