package device

import (
	"fmt"
	"sync"
)

type MockModbus struct {
	mu            sync.RWMutex
	unitId        uint8
	registers     map[uint16]uint16
	shouldError   bool
	errorMessage  string
	recordingMode bool
	recordingFile string
}

// NewMockModbus creates a new mock Modbus client
func NewMockModbus() *MockModbus {
	return &MockModbus{
		registers: make(map[uint16]uint16),
	}
}

func (m *MockModbus) SetUnitId(id uint8) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.shouldError {
		return fmt.Errorf("SetUnitId: %s", m.errorMessage)
	}

	m.unitId = id
	return nil
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

// SetRegister sets a register value for testing
func (m *MockModbus) SetRegister(address uint16, value uint16) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.registers[address] = value
}

// SetRegisters sets multiple register values
func (m *MockModbus) SetRegisters(startAddr uint16, values []uint16) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, val := range values {
		m.registers[startAddr+uint16(i)] = val
	}
}

// SimulateError configures the mock to return errors
func (m *MockModbus) SimulateError(shouldError bool, message string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shouldError = shouldError
	m.errorMessage = message
}
