// File: internal/device/pxu_test.go
package device

import (
	"testing"
	"time"
)

func TestNewPxu(t *testing.T) {
	tests := []struct {
		name        string
		unitId      uint8
		client      Modbus
		timeout     time.Duration
		retries     int
		expectError bool
	}{
		{
			name:        "valid parameters",
			unitId:      1,
			client:      NewMockModbus(),
			timeout:     time.Second,
			retries:     3,
			expectError: false,
		},
		{
			name:        "nil client",
			unitId:      1,
			client:      nil,
			timeout:     time.Second,
			retries:     3,
			expectError: true,
		},
		{
			name:        "zero timeout uses default",
			unitId:      1,
			client:      NewMockModbus(),
			timeout:     0,
			retries:     3,
			expectError: false,
		},
		{
			name:        "zero retries uses default",
			unitId:      1,
			client:      NewMockModbus(),
			timeout:     time.Second,
			retries:     0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pxu, err := NewPxu(tt.unitId, tt.client, tt.timeout, tt.retries)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if pxu == nil {
				t.Error("expected non-nil PXU instance")
				return
			}

			// Check defaults were applied
			if tt.timeout == 0 && pxu.timeout != DefaultTimeout {
				t.Errorf("expected default timeout %v, got %v", DefaultTimeout, pxu.timeout)
			}

			if tt.retries == 0 && pxu.retries != DefaultRetries {
				t.Errorf("expected default retries %d, got %d", DefaultRetries, pxu.retries)
			}
		})
	}
}

func TestPxu_ReadDeviceStats(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func(*MockModbus)
		expectError  bool
		validateFunc func(*testing.T, *Stats)
	}{
		{
			name: "successful read with celsius",
			setupMock: func(mock *MockModbus) {
				// Set up register values for a successful read
				registers := make([]uint16, StatsRegCount)
				registers[StatsRegPV] = 250                   // 25.0°C
				registers[StatsRegSP] = 300                   // 30.0°C
				registers[StatsRegLED] = LEDCelsius | LEDOut1 // Celsius + Out1 active
				registers[StatsRegRS] = uint16(RunStatusRun)

				mock.SetRegisters(0, registers)
			},
			expectError: false,
			validateFunc: func(t *testing.T, stats *Stats) {
				if stats.Pv != 25.0 {
					t.Errorf("expected PV 25.0, got %f", stats.Pv)
				}
				if stats.Sp != 30.0 {
					t.Errorf("expected SP 30.0, got %f", stats.Sp)
				}
				if stats.VUnit != "C" {
					t.Errorf("expected unit 'C', got '%s'", stats.VUnit)
				}
				if !stats.Out1 {
					t.Error("expected Out1 to be true")
				}
				if stats.Out2 {
					t.Error("expected Out2 to be false")
				}
				if stats.RS != RunStatusRun {
					t.Errorf("expected run status RUN, got %v", stats.RS)
				}
			},
		},
		{
			name: "successful read with fahrenheit",
			setupMock: func(mock *MockModbus) {
				registers := make([]uint16, StatsRegCount)
				registers[StatsRegPV] = 770                      // 77.0°F
				registers[StatsRegSP] = 860                      // 86.0°F
				registers[StatsRegLED] = LEDFahrenheit | LEDOut2 // Fahrenheit + Out2 active

				mock.SetRegisters(0, registers)
			},
			expectError: false,
			validateFunc: func(t *testing.T, stats *Stats) {
				if stats.VUnit != "F" {
					t.Errorf("expected unit 'F', got '%s'", stats.VUnit)
				}
				if stats.Out1 {
					t.Error("expected Out1 to be false")
				}
				if !stats.Out2 {
					t.Error("expected Out2 to be true")
				}
			},
		},

		// this error happens when the new device is created.  we can argue that this should happen in each function
		// instead....

		//{
		//	name: "modbus read error",
		//	setupMock: func(mock *MockModbus) {
		//		mock.SimulateError(true, "connection timeout")
		//	},
		//	expectError: true,
		//},
		{
			name: "insufficient registers",
			setupMock: func(mock *MockModbus) {
				// Only set up 10 registers instead of 30
				registers := make([]uint16, 10)
				mock.SetRegisters(0, registers)
			},
			expectError: true,
		},
		{
			name: "invalid temperature unit",
			setupMock: func(mock *MockModbus) {
				registers := make([]uint16, StatsRegCount)
				registers[StatsRegLED] = 0 // No temperature unit bits set

				mock.SetRegisters(0, registers)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockModbus()
			tt.setupMock(mock)

			pxu, err := NewPxu(1, mock, time.Second, 1)
			if err != nil {
				t.Fatalf("failed to create PXU: %v", err)
			}

			stats, err := pxu.ReadStats()

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if stats == nil {
				t.Error("expected non-nil stats")
				return
			}

			if tt.validateFunc != nil {
				tt.validateFunc(t, stats)
			}
		})
	}
}

func TestPxu_ReadDeviceInfo(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func(*MockModbus)
		expectError  bool
		validateFunc func(*testing.T, *Info)
	}{
		{
			name: "successful read",
			setupMock: func(mock *MockModbus) {
				// Simulate "PXU123" as hex values + firmware version
				registers := []uint16{
					0x5058, // "PX"
					0x5531, // "U1"
					0x3233, // "23"
					0x0000, // padding
					0x0000, // padding
					0x0000, // padding
					123,    // firmware 1.23
				}
				mock.SetRegisters(InfoRegStart, registers)
			},
			expectError: false,
			validateFunc: func(t *testing.T, info *Info) {
				if info.Model == "" {
					t.Error("expected non-empty model")
				}
				if info.Firmware != "1.23" {
					t.Errorf("expected firmware '1.23', got '%s'", info.Firmware)
				}
			},
		},
		//{
		//	name: "modbus read error",
		//	setupMock: func(mock *MockModbus) {
		//		mock.SimulateError(true, "device not responding")
		//	},
		//	expectError: true,
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockModbus()
			tt.setupMock(mock)

			pxu, err := NewPxu(1, mock, time.Second, 1)
			if err != nil {
				t.Fatalf("failed to create PXU: %v", err)
			}

			info, err := pxu.ReadInfo()

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if info == nil {
				t.Error("expected non-nil info")
				return
			}

			if tt.validateFunc != nil {
				tt.validateFunc(t, info)
			}
		})
	}
}

//func TestPxu_RetryLogic(t *testing.T) {
//	mock := NewMockModbus()
//
//	// Set up mock to fail first 2 attempts, succeed on 3rd
//	callCount := 0
//	originalRead := mock.ReadRegisters
//	mock.ReadRegisters = func(address, quantity uint16) ([]uint16, error) {
//		callCount++
//		if callCount < 3 {
//			return nil, fmt.Errorf("simulated failure %d", callCount)
//		}
//		// Reset error state and call original
//		mock.SimulateError(false, "")
//		return originalRead(address, quantity)
//	}
//
//	registers := make([]uint16, StatsRegCount)
//	registers[StatsRegLED] = LEDCelsius
//	mock.SetRegisters(0, registers)
//
//	device, err := NewPxu(1, mock, time.Second, 3)
//	if err != nil {
//		t.Fatalf("failed to create PXU: %v", err)
//	}
//
//	stats, err := device.ReadStats()
//	if err != nil {
//		t.Errorf("expected success after retries, got error: %v", err)
//	}
//
//	if stats == nil {
//		t.Error("expected non-nil stats after successful retry")
//	}
//
//	if callCount != 3 {
//		t.Errorf("expected 3 calls (2 failures + 1 success), got %d", callCount)
//	}
//}

func TestPxu_ReadProfile(t *testing.T) {

	tests := []struct {
		name          string
		profileNumber uint8
		setupMock     func(*MockModbus)
		expectError   bool
		validateFunc  func(*testing.T, *Info)
	}{
		{
			name:          "successful read",
			profileNumber: 0,
			setupMock: func(mock *MockModbus) {
				// Simulate "PXU123" as hex values + firmware version
				registers := []uint16{
					0x5058, // "PX"
					0x5531, // "U1"
					0x3233, // "23"
					0x0000, // padding
					0x0000, // padding
					0x0000, // padding
					123,    // firmware 1.23
				}
				mock.SetRegisters(InfoRegStart, registers)
			},
			expectError: false,
			validateFunc: func(t *testing.T, info *Info) {
				if info.Model == "" {
					t.Error("expected non-empty model")
				}
				if info.Firmware != "1.23" {
					t.Errorf("expected firmware '1.23', got '%s'", info.Firmware)
				}
			},
		},
		{
			name:          "profile number exceeded",
			profileNumber: 17,
			expectError:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockModbus()
			if tt.setupMock != nil {
				tt.setupMock(mock)
			}
			pxu, err := NewPxu(1, mock, time.Second, 1)
			if err != nil {
				t.Fatalf("failed to create PXU: %v", err)
			}

			err = pxu.ReadProfile(tt.profileNumber)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			//if info == nil {
			//	t.Error("expected non-nil info")
			//	return
			//}
			//
			//if tt.validateFunc != nil {
			//	tt.validateFunc(t, info)
			//}
		})
	}
}
