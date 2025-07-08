package device

import (
	"fmt"
	"time"
)

type Driver struct {
	cfg        Configuration
	unitId     uint8
	controller PidController
}

func NewDriver(unitId uint8, cfg Configuration) (*Driver, error) {
	client, err := NewModbusHandler(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate modbus handler for unit %d: %v", unitId, err)
	}

	pxu, err := NewPxu(unitId, client, 3*time.Second, 30)
	if err != nil {
		return nil, fmt.Errorf("failed to create PID controller for unit %d: %v", unitId, err)
	}
	return &Driver{cfg: cfg, unitId: unitId, controller: pxu}, nil
}

func (d *Driver) Close() error {
	return d.controller.Close()
}
