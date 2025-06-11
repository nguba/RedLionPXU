// File: internal/device/pxu_test.go
package device

import (
	"reflect"
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
				registers[RegPV] = 250                   // 25.0°C
				registers[RegSP] = 300                   // 30.0°C
				registers[RegLED] = LEDCelsius | LEDOut1 // Celsius + Out1 active
				registers[RegControllerStatus] = uint16(RunStatusRun)

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
				registers[RegPV] = 770                      // 77.0°F
				registers[RegSP] = 860                      // 86.0°F
				registers[RegLED] = LEDFahrenheit | LEDOut2 // Fahrenheit + Out2 active

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
				registers[RegLED] = 0 // No temperature unit bits set

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

func TestPxu_ReadInfo(t *testing.T) {

	// TODO add tests for reading profile segments correctly
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
				mock.SetRegisters(RegInfoStart, registers)
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

func TestPxu_ReadProfile(t *testing.T) {
	tests := []struct {
		name          string
		profileNumber uint16
		setupMock     func(*MockModbus)
		expectError   bool
		validateFunc  func(*testing.T, *Profile)
	}{
		{
			name: "successful read",
			setupMock: func(mock *MockModbus) {
				registers := []uint16{
					4, // 4 segments
				}
				mock.SetRegisters(RegNumSegmentsStart, registers)

				// Simulate "PXU123" as hex values + firmware version
				registers = []uint16{
					250,  // 25.0C
					7200, // 12 hr
					305,  // 30.5C
					3600, // 6 hr
					620,  // 62.0C
					7200, // 12 hr
					720,  // 72C
					9999, // 999.9 minutes
				}

				mock.SetRegisters(RegProfSegmentStart, registers)
			},
			expectError: false,
			validateFunc: func(t *testing.T, profile *Profile) {
				if profile.Id != 0 {
					t.Errorf("expected profile with id '%d', got '%d'", profile.Id, 0)
				}
				if profile.SegCount == 0 || profile.Segments == nil {
					t.Error("expected non-empty segments")
				}
				if len(profile.Segments) != int(profile.SegCount) {
					t.Errorf("profile segment count mismatch, want '%d', got '%d'", profile.SegCount, len(profile.Segments))
				}
				var expected []Segment
				expected = append(expected, Segment{
					Id: 0,
					Sp: 250,
					T:  7200,
				})
				expected = append(expected, Segment{
					Id: 1,
					Sp: 305,
					T:  3600,
				})
				expected = append(expected, Segment{
					Id: 2,
					Sp: 620,
					T:  7200,
				})
				expected = append(expected, Segment{
					Id: 3,
					Sp: 720,
					T:  9999,
				})

				if !reflect.DeepEqual(profile.Segments, expected) {
					t.Errorf("profile segments mismatch, want '%v', got '%v'", expected, profile.Segments)
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

			profile, err := pxu.ReadProfile(tt.profileNumber)

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

			if profile == nil {
				t.Error("expected non-nil profile")
				return
			}

			if tt.validateFunc != nil {
				tt.validateFunc(t, profile)
			}
		})
	}
}

func TestPxu_Start(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func(*MockModbus)
		expectError  bool
		validateFunc func(*testing.T, *Profile)
	}{
		{
			name: "start temperature control",
			setupMock: func(mock *MockModbus) {
				_ = mock.SetRegister(RegControllerStatus, RsStop)
			},
			expectError: false,
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

			val, err := mock.ReadRegister(RegControllerStatus)

			if err != nil {
				t.Fatalf("failed to read register: %v", err)
			}
			if val != RsStop {
				t.Errorf("expected value '%d', got '%d'", RsStop, val)
			}

			err = pxu.Start()

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
		})
	}
}
