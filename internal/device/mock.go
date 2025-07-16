package device

import (
	"fmt"
	"sync"
)

type MockModbus struct {
	mu            sync.RWMutex
	unitId        UnitId
	registers     map[uint16]uint16
	shouldError   bool
	errorMessage  string
	recordingMode bool
	recordingFile string
}

// NewMockModbus creates a new mock Modbus client impersonating the RedLion PXU.  A new Pxu can be instantiated
// with this test double as the client when no device is plugged in.
func NewMockModbus() *MockModbus {
	return &MockModbus{
		registers: make(map[uint16]uint16),
	}
}

func (m *MockModbus) SetUnitId(id UnitId) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.shouldError {
		return fmt.Errorf("SetUnitId: %s", m.errorMessage)
	}

	m.unitId = id
	return nil
}

func (m *MockModbus) ReadRegister(address uint16) (uint16, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.shouldError {
		return ErrVal, fmt.Errorf("ReadRegister: %s", m.errorMessage)
	}
	if val, exists := m.registers[address]; exists {
		return val, nil
	}
	return ErrVal, nil
}

func (m *MockModbus) ReadRegisters(address, quantity uint16) ([]uint16, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.shouldError {
		return nil, fmt.Errorf("ReadRegisters: %s", m.errorMessage)
	}

	// Build response from stored registers
	response := make([]uint16, quantity)
	for i := uint16(0); i < quantity; i++ {
		if val, exists := m.registers[address+i]; exists {
			response[i] = val
		} else {
			response[i] = 0 // Default value for unset registers
		}
	}
	return response, nil
}

func (m *MockModbus) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return nil
}

func (m *MockModbus) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	clear(m.registers)
}

// SetRegister sets a register value for testing
func (m *MockModbus) SetRegister(address uint16, value uint16) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.registers[address] = value

	return nil
}

// SetRegisters sets multiple register values
func (m *MockModbus) SetRegisters(startAddr uint16, values []uint16) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, val := range values {
		m.registers[startAddr+uint16(i)] = val
	}

	return nil
}

// SimulateError configures the mock to return errors
func (m *MockModbus) SimulateError(shouldError bool, message string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shouldError = shouldError
	m.errorMessage = message
}

func (m *MockModbus) GetStatsRegister() []uint16 {
	registers := make([]uint16, StatsRegCount)
	registers[RegPV] = 255                   // 25.5°C
	registers[RegSP] = 304                   // 30.4°C
	registers[RegLED] = LEDCelsius | LEDOut1 // Celsius + Out1 active
	registers[RegControllerStatus] = uint16(RunStatusRun)

	return registers
}
