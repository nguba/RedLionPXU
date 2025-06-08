package device

import (
	"fmt"
	"log"
	"time"
)

type Pxu struct {
	client  Modbus
	timeout time.Duration
	retries int
	unitId  uint8
}

func NewPxu(unitId uint8, client Modbus, timeout time.Duration, retries int) (*Pxu, error) {

	if client == nil {
		return nil, fmt.Errorf("modbus client cannot be nil")
	}

	if timeout == 0 {
		timeout = DefaultTimeout
	}

	if retries == 0 {
		retries = DefaultRetries
	}

	if err := client.SetUnitId(unitId); err != nil {
		return nil, fmt.Errorf("failed to set unit ID %d: %w", unitId, err)
	}

	log.Printf("connected to unit %d", unitId)

	controller := &Pxu{
		client:  client,
		timeout: timeout,
		retries: retries,
		unitId:  unitId,
	}
	return controller, nil
}

func (p *Pxu) readRegistersWithRetry(addr, count uint16) ([]uint16, error) {
	var lastErr error

	for attempt := 0; attempt <= p.retries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			backoff := time.Duration(attempt) * 100 * time.Millisecond
			time.Sleep(backoff)
		}

		regs, err := p.client.ReadRegisters(addr, count)
		if err == nil {
			return regs, nil
		}
		lastErr = err
	}

	return nil, fmt.Errorf("failed after %d retries: %w", p.retries, lastErr)
}

func (p *Pxu) Close() error {
	if p.client != nil {
		return nil
	}
	return p.client.Close()
}

func (p *Pxu) ReadStats() (*Stats, error) {
	const totalRegisters = 30

	regs, err := p.readRegistersWithRetry(0, totalRegisters)
	if err != nil {
		return nil, fmt.Errorf("failed reading registers from unit %d: %w", p.unitId, err)
	}

	// ensure we received all the requested registers
	if len(regs) < totalRegisters {
		return nil, fmt.Errorf("insufficient registers received: expected %d, got %d", totalRegisters, len(regs))
	}

	return NewStats(regs)
}

func (p *Pxu) ReadInfo() (*Info, error) {
	regs, err := p.readRegistersWithRetry(InfoRegStart, InfoRegCount)
	if err != nil {
		return nil, fmt.Errorf("failed reading info registers from unit %d: %w", p.unitId, err)
	}

	if len(regs) < InfoRegCount {
		return nil, fmt.Errorf("insufficient registers received: expected %d, got %d", InfoRegCount, len(regs))
	}

	return NewInfo(regs)
}

func (p *Pxu) ReadProfile(number uint8) error {
	if number > 16 {
		return fmt.Errorf("invalid profile selected: %d", number)
	}

	profile := Profile{Num: number}

	// read the number of segments this profile spans
	segReg, err := p.readRegistersWithRetry(ProfNumSegmentsRegStart, 1)
	if err != nil {
		return fmt.Errorf("failed reading profile segment count from unit %d: %w", p.unitId, err)
	}
	// this is how many segments are configured in this profile.
	profile.SegCount = segReg[0]

	reg, err := p.readRegistersWithRetry(ProfSegmentRegStart, profile.SegCount*2)
	if err != nil {
		return fmt.Errorf("failed reading profile from unit %d: %w", p.unitId, err)
	}

	var Sp uint16 = 0
	for i := uint16(0); i < profile.SegCount; i++ {
		seg := Segment{
			Num: uint8(i),
			Sp:  reg[Sp],
			T:   reg[Sp+1],
		}
		profile.Segments = append(profile.Segments, seg)
		Sp += 2
	}

	fmt.Printf("profile: %+v\n", profile)

	return nil
}
