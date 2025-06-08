package device

import (
	"fmt"
	"strings"
	"time"
)

type Stats struct {
	Pv     float64   `json:"pv"`
	Sp     float64   `json:"sp"`
	Out1   bool      `json:"out1"`
	Out2   bool      `json:"out2"`
	At     bool      `json:"at"`
	TP     float64   `json:"tp"`
	TI     uint16    `json:"ti"`
	TD     uint16    `json:"td"`
	TGroup uint16    `json:"tgroup"`
	RS     RunStatus `json:"rs"`
	VUnit  string    `json:"vunit"`
	PC     uint16    `json:"pc"`
	PS     uint16    `json:"ps"`
	PSR    float64   `json:"psr"`
}

type Info struct {
	Model    string `json:"model"`
	Firmware string `json:"firmware"`
}

func (info Info) String() string {
	return fmt.Sprintf("Model: %s, Firmware: %s", strings.Trim(info.Model, " "), info.Firmware)
}

// RunStatus represents device run status
type RunStatus uint16

const (
	RunStatusStop RunStatus = iota
	RunStatusRun
	RunStatusEnd
	RunStatusPause
	RunStatusAdvanceProfile
)

func (s RunStatus) String() string {
	switch s {
	case RunStatusStop:
		return "STOP"
	case RunStatusRun:
		return "RUN"
	case RunStatusEnd:
		return "END"
	case RunStatusPause:
		return "PAUSE"
	case RunStatusAdvanceProfile:
		return "ADVANCE PROFILE"
	default:
		return fmt.Sprintf("UNKNOWN (%d)", s)
	}
}

// Configuration stores the settings for communication with serial devices.
// Example URL:  rtu://COM3 (windows), rtu:///dev/ttyUSB0 (Linux).
type Configuration struct {
	URL      string
	Speed    uint
	DataBits uint
	Parity   string
	Timeout  time.Duration
}

type Segment struct {
	Num uint8
	Sp  uint16
	T   uint16
}

type Profile struct {
	Num      uint8
	SegCount uint16
	Segments []Segment
}
