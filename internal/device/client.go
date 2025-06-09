package device

import (
	"fmt"
	"github.com/simonvetter/modbus"
	"log"
)

// ModbusHandler wraps the modbus client
type ModbusHandler struct {
	modbus *modbus.ModbusClient
}

func (c *ModbusHandler) SetUnitId(id uint8) error {
	if c.modbus == nil {
		return fmt.Errorf("modbus client is nil")
	}
	return c.modbus.SetUnitId(id)
}

func (c *ModbusHandler) ReadRegisters(address, quantity uint16) ([]uint16, error) {
	if c.modbus == nil {
		return nil, fmt.Errorf("modbus client is nil")
	}
	regs, err := c.modbus.ReadRegisters(address, quantity, modbus.HOLDING_REGISTER)
	if err != nil {
		log.Printf("failed to read registers addr=%d, qty=%d: %v", address, quantity, err)
	}
	return regs, err
}

func (c *ModbusHandler) Close() error {
	if c.modbus == nil {
		return nil
	}
	return c.modbus.Close()
}

// NewModbusHandler creates a new PXU client with the given configuration
func NewModbusHandler(cfg *Configuration) (*ModbusHandler, error) {

	modbusConfig := &modbus.ClientConfiguration{
		URL:      cfg.URL,
		Speed:    cfg.Speed,    // 38400
		DataBits: cfg.DataBits, // 8
		Parity:   modbus.PARITY_NONE,
		Timeout:  cfg.Timeout, //  500 * time.Millisecond
		Logger:   log.Default(),
	}
	client, err := modbus.NewClient(modbusConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating modbus client: %w", err)
	}

	err = client.Open() // needed for communicating with this device
	if err != nil {
		return nil, fmt.Errorf("error opening modbus connection: %w", err)
	}
	return &ModbusHandler{modbus: client}, nil
}
