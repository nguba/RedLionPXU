package main

import (
	"github.com/nguba/RedLionPXU/internal/device"
	"log"
	"time"
)

const unit = device.UnitId(5)

// DefaultConfiguration returns a default configuration for COM3
func DefaultConfiguration() *device.Configuration {
	return &device.Configuration{
		URL:      "rtu://COM3",
		Speed:    device.DefaultSpeed,
		DataBits: 8,
		Parity:   "none",
		Timeout:  500 * time.Millisecond,
	}
}

// NewDefaultPxu creates a PXU with default settings
func NewDefaultPxu(unitId device.UnitId) (*device.Pxu, error) {
	cfg := DefaultConfiguration()
	client, err := device.NewModbusDevice(cfg)
	if err != nil {
		return nil, err
	}
	return device.NewPxu(unitId, client, device.DefaultTimeout, device.DefaultRetries)
}

func main() {
	log.Printf("starting server for unit %d", unit)

	pxu, err := NewDefaultPxu(unit)
	if err != nil {
		log.Fatal(err)
	}
	defer func(pxu *device.Pxu) {
		_ = pxu.Close()
	}(pxu)
}
