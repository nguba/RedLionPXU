package device

import (
	"testing"
)

func TestMockModbus_BasicOperations(t *testing.T) {
	mock := NewMockModbus()

	// Test SetUnitId
	err := mock.SetUnitId(5)
	if err != nil {
		t.Errorf("unexpected error setting unit ID: %v", err)
	}

	// Test setting and reading registers
	mock.SetRegister(100, 1234)
	mock.SetRegister(101, 5678)

	regs, err := mock.ReadRegisters(100, 2)
	if err != nil {
		t.Errorf("unexpected error reading registers: %v", err)
	}

	if len(regs) != 2 {
		t.Errorf("expected 2 registers, got %d", len(regs))
	}

	if regs[0] != 1234 || regs[1] != 5678 {
		t.Errorf("expected [1234, 5678], got %v", regs)
	}
}

func TestMockModbus_ErrorSimulation(t *testing.T) {
	mock := NewMockModbus()

	// Configure mock to return errors
	mock.SimulateError(true, "test error")

	err := mock.SetUnitId(1)
	if err == nil {
		t.Error("expected error from SetUnitId but got none")
	}

	_, err = mock.ReadRegisters(0, 10)
	if err == nil {
		t.Error("expected error from ReadRegisters but got none")
	}
}
