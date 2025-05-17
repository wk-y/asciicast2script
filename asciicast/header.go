// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package asciicast

import (
	"encoding/json"
	"fmt"
)

type UnsupportedVersionError struct {
	Version int
}

func (u UnsupportedVersionError) Error() string {
	return fmt.Sprintf("asciicast v%d is unsupported", u.Version)
}

type Header interface {
	Version() int
	Width() int
	Height() int
	Term() (term string, ok bool) // ex. "xterm-256color"
	Timestamp() (timestamp int64, ok bool)
	Duration() (duration float64, ok bool)
	Command() (command string, ok bool)
	Title() (title string, ok bool)
	IdleTimeLimit() (idleTimeLimit int, ok bool)
	Env() map[string]string
	Theme() map[string]string
	RelativeTime() bool // interpret time as relative
}

func DecodeHeader(rawHeader []byte) (Header, error) {
	versionExtractor := struct {
		Version int `json:"version"`
	}{}

	if err := json.Unmarshal(rawHeader, &versionExtractor); err != nil {
		return nil, err
	}

	switch versionExtractor.Version {
	case 2:
		var h HeaderV2
		err := json.Unmarshal(rawHeader, &h)
		return HeaderV2Iface{Header: h}, err
	case 3:
		var h HeaderV3
		err := json.Unmarshal(rawHeader, &h)
		return HeaderV3Iface{Header: h}, err
	default:
		return nil, UnsupportedVersionError{Version: versionExtractor.Version}
	}
}

type HeaderV2 struct {
	Version       int               `json:"version"`
	Width         int               `json:"width"`
	Height        int               `json:"height"`
	Timestamp     *int64            `json:"timestamp"`
	Duration      *float64          `json:"duration"`
	Command       *string           `json:"command"`
	Title         *string           `json:"title"`
	IdleTimeLimit *int              `json:"idle_time_limit"`
	Env           map[string]string `json:"env"`
	Theme         map[string]string `json:"theme"`
}

// Wrapper to avoid field/method collisions
type HeaderV2Iface struct {
	Header HeaderV2
}

var _ Header = HeaderV2Iface{}

func (h HeaderV2Iface) Version() int {
	return h.Header.Version
}

func (h HeaderV2Iface) Width() int {
	return h.Header.Width
}

func (h HeaderV2Iface) Height() int {
	return h.Header.Height
}

func (h HeaderV2Iface) Term() (term string, ok bool) {
	term, ok = h.Header.Env["TERM"]
	return
}

func (h HeaderV2Iface) Timestamp() (timestamp int64, ok bool) {
	if h.Header.Timestamp == nil {
		return 0, false
	}
	return *h.Header.Timestamp, true
}

func (h HeaderV2Iface) Duration() (duration float64, ok bool) {
	if h.Header.Duration == nil {
		return 0, false
	}
	return *h.Header.Duration, true
}

func (h HeaderV2Iface) Command() (command string, ok bool) {
	if h.Header.Command == nil {
		return "", false
	}
	return *h.Header.Command, true
}

func (h HeaderV2Iface) Title() (title string, ok bool) {
	if h.Header.Title == nil {
		return "", false
	}
	return *h.Header.Title, true
}

func (h HeaderV2Iface) IdleTimeLimit() (idleTimeLimit int, ok bool) {
	if h.Header.IdleTimeLimit == nil {
		return 0, false
	}
	return *h.Header.IdleTimeLimit, true
}

func (h HeaderV2Iface) Env() map[string]string {
	return h.Header.Env
}

func (h HeaderV2Iface) Theme() map[string]string {
	return h.Header.Theme
}

func (h HeaderV2Iface) RelativeTime() bool {
	return false
}

type TermInfo struct {
	Cols    int               `json:"cols"`
	Rows    int               `json:"rows"`
	Type    *string           `json:"type"`
	Version *string           `json:"version"`
	Theme   map[string]string `json:"theme"`
}

type HeaderV3 struct {
	Version       int               `json:"version"`
	Term          TermInfo          `json:"term"`
	Timestamp     *int64            `json:"timestamp"`
	Duration      *float64          `json:"duration"`
	Command       *string           `json:"command"`
	Title         *string           `json:"title"`
	IdleTimeLimit *int              `json:"idle_time_limit"`
	Env           map[string]string `json:"env"`
}

// Wrapper to avoid field/method collisions
type HeaderV3Iface struct {
	Header HeaderV3
}

var _ Header = HeaderV3Iface{}

func (h HeaderV3Iface) Version() int {
	return h.Header.Version
}

func (h HeaderV3Iface) Width() int {
	return h.Header.Term.Cols
}

func (h HeaderV3Iface) Height() int {
	return h.Header.Term.Rows
}

func (h HeaderV3Iface) Term() (term string, ok bool) {
	if h.Header.Term.Type == nil {
		return "", false
	}
	return *h.Header.Term.Type, true
}

func (h HeaderV3Iface) Timestamp() (timestamp int64, ok bool) {
	if h.Header.Timestamp == nil {
		return 0, false
	}
	return *h.Header.Timestamp, true
}

func (h HeaderV3Iface) Duration() (duration float64, ok bool) {
	if h.Header.Duration == nil {
		return 0, false
	}
	return *h.Header.Duration, true
}

func (h HeaderV3Iface) Command() (command string, ok bool) {
	if h.Header.Command == nil {
		return "", false
	}
	return *h.Header.Command, true
}

func (h HeaderV3Iface) Title() (title string, ok bool) {
	if h.Header.Title == nil {
		return "", false
	}
	return *h.Header.Title, true
}

func (h HeaderV3Iface) IdleTimeLimit() (idleTimeLimit int, ok bool) {
	if h.Header.IdleTimeLimit == nil {
		return 0, false
	}
	return *h.Header.IdleTimeLimit, true
}

func (h HeaderV3Iface) Env() map[string]string {
	return h.Header.Env
}

func (h HeaderV3Iface) Theme() map[string]string {
	return h.Header.Term.Theme
}

func (h HeaderV3Iface) RelativeTime() bool {
	return true
}
