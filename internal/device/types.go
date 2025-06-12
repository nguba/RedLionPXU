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

func (s Stats) String() string {
	return fmt.Sprintf(
		"PV:%.1f%s SP:%.1f%s | Out1:%t Out2:%t AT:%t | TP:%.1f TI:%d TD:%d TGroup:%d | RS:%s | PC:%d PS:%d PSR:%.1f",
		s.Pv, s.VUnit,
		s.Sp, s.VUnit,
		s.Out1, s.Out2, s.At,
		s.TP, s.TI, s.TD, s.TGroup,
		s.RS,
		s.PC, s.PS, s.PSR,
	)
}

type Info struct {
	Model    string `json:"model"`
	Firmware string `json:"firmware"`
}

func (info Info) String() string {
	return fmt.Sprintf("Model: %s, Firmware: %s", strings.Trim(info.Model, " "), info.Firmware)
}

// RunStatus shows the state of the device: RunStatusStop, RunStatusRun, RunStatusEnd, RunStatusPause, RunStatusAdvanceProfile
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
	Id uint8
	Sp float64
	T  float64
}

func (s Segment) String() string {
	return fmt.Sprintf("Id: %d, Sp: %.1f, T: %.1f", s.Id, s.Sp, s.T)
}

type Profile struct {
	Id       uint16
	SegCount uint16
	link     uint16
	Segments []Segment
	repeat   uint16
}

func (p Profile) String() string {
	var linkVal string
	if p.link == LinkEnd {
		linkVal = "END"
	} else if p.link == LinkStop {
		linkVal = "STOP"
	} else {
		linkVal = fmt.Sprintf("PROFILE %d", p.link)
	}
	return fmt.Sprintf("Id: %d, SegCount: %d, Link: %s, Repeat: %d, Segments: %+v",
		p.Id, p.SegCount, linkVal, p.repeat, p.Segments)
}

func NewProfile(id uint16, segmentCount, linkProfile, repeatCycle uint16) *Profile {
	profile := Profile{Id: id}
	profile.SegCount = segmentCount // configured active segments
	profile.link = linkProfile      // profile to continue with, END or STOP
	profile.repeat = repeatCycle    // whether this profile repeats and how often (0 = no repeat)
	return &profile
}
