package device

import (
	"errors"
	"fmt"
	"github.com/simonvetter/modbus"
	"log"
)

// ModbusHandler wraps the modbus client
type ModbusHandler struct {
	modbus *modbus.ModbusClient
}

func (c *ModbusHandler) SetUnitId(id UnitId) error {
	if c.modbus == nil {
		return fmt.Errorf("modbus client is nil")
	}
	return c.modbus.SetUnitId(uint8(id))
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

func (c *ModbusHandler) SetRegister(address, value uint16) error {
	if c.modbus == nil {
		return fmt.Errorf("modbus client is nil")
	}

	err := c.modbus.WriteRegister(address, value)
	if err != nil {
		return fmt.Errorf("error writing registersers addr=%d, value:%d: %v", address, value, err)
	}

	return nil
}

func (m *ModbusHandler) SetRegisters(startAddr uint16, values []uint16) error {
	return errors.New("not implemented")
}

func (c *ModbusHandler) ReadRegister(address uint16) (uint16, error) {
	if c.modbus == nil {
		return ErrVal, fmt.Errorf("modbus client is nil")
	}

	val, err := c.modbus.ReadRegister(address, modbus.HOLDING_REGISTER)
	if err != nil {
		return ErrVal, fmt.Errorf("error reading registerser addr=%d: %w", address, err)
	}

	return val, nil
}
