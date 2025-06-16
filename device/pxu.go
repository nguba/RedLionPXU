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

	log.Printf("connecting to unit %d", unitId)
	if err := client.SetUnitId(unitId); err != nil {
		return nil, fmt.Errorf("failed to set unit ID %d: %w", unitId, err)
	}

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
	regs, err := p.readRegistersWithRetry(RegInfoStart, InfoRegCount)
	if err != nil {
		return nil, fmt.Errorf("failed reading registers from unit %d: %w", p.unitId, err)
	}

	if len(regs) < InfoRegCount {
		return nil, fmt.Errorf("insufficient registers received: expected %d, got %d", InfoRegCount, len(regs))
	}

	return NewInfo(regs)
}

func (p *Pxu) ReadProfile(id uint16) (*Profile, error) {
	if id > 16 {
		return nil, fmt.Errorf("invalid profile id selected: %d", id)
	}

	// read the number of segments this profile spans
	segmentCount, err := p.readRegistersWithRetry(RegNumSegments+id, 1)
	if err != nil {
		return nil, fmt.Errorf("failed reading profile segment count from unit %d: %w", p.unitId, err)
	}

	// read whether the profile stops, ends or continues with another one
	linkProfile, err := p.readRegistersWithRetry(RegProfLink+id, 1)
	if err != nil {
		return nil, fmt.Errorf("failed reading linked profile from unit %d: %w", p.unitId, err)
	}

	// read how often the profile repeats
	repeatCycle, err := p.readRegistersWithRetry(RegProfCycleRepeat+id, 1)
	if err != nil {
		return nil, fmt.Errorf("failed reading profile cycle count from unit %d: %w", p.unitId, err)
	}

	sc := segmentCount[0] + 1 // count of zero actually means one segment only
	profile := NewProfile(id, sc, linkProfile[0], repeatCycle[0])

	start := id*32 + RegProfSegmentStart
	count := profile.SegCount * 2
	regs, err := p.readRegistersWithRetry(start, count)
	if err != nil {
		return nil, fmt.Errorf("failed reading profile from unit %d: %w", p.unitId, err)
	}

	fillProfile(profile, regs)
	return profile, nil
}

func fillProfile(profile *Profile, regs []uint16) {
	// setpoint -> even idx, time -> odd idx
	for i := uint16(0); i < profile.SegCount; i++ {
		p := i * 2
		seg := Segment{
			Id: uint8(i),
			Sp: toFloat(regs[p]),
			T:  toFloat(regs[p+1]),
		}
		profile.Segments = append(profile.Segments, seg)
	}
}

func (p *Pxu) UpdateSetpoint(value float64) error {
	err := p.client.SetRegister(RegSP, toUint16(value))
	if err != nil {
		return fmt.Errorf("failed to update sp to %.1f: %w", value, err)
	}
	log.Printf("updated sp to %.1f", value)
	return nil
}

func (p *Pxu) UpdateControllerStatus(value uint16) error {
	err := p.client.SetRegister(RegControllerStatus, value)
	if err != nil {
		return fmt.Errorf("failed to update controller status on unit %d: %w", p.unitId, err)
	}
	return nil
}

func (p *Pxu) Stop() error {
	if err := p.UpdateControllerStatus(RsStop); err != nil {
		return fmt.Errorf("failed to stop unit %d: %w", p.unitId, err)
	}
	log.Printf("stopped unit %d", p.unitId)
	return nil
}

func (p *Pxu) Start() error {
	if err := p.UpdateControllerStatus(RsStart); err != nil {
		return fmt.Errorf("failed to start unit %d: %w", p.unitId, err)
	}
	log.Printf("started unit %d", p.unitId)
	return nil
}
