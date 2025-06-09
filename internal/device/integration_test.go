package device

import (
	"fmt"
	"testing"
	"time"
)

// TestIntegration_DeviceStatsWorkflow tests the complete workflow
func TestIntegration_StatsWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Set up mock with realistic data
	mock := NewMockModbus()

	// Configure device stats registers
	statsRegs := make([]uint16, StatsRegCount)
	statsRegs[RegPV] = 235   // 23.5°C
	statsRegs[RegSP] = 250   // 25.0°C
	statsRegs[RegTP] = 180   // 18.0°C
	statsRegs[RegTI] = 120   // TI value
	statsRegs[RegTD] = 30    // TD value
	statsRegs[RegTGroup] = 1 // Temperature group
	statsRegs[RegRS] = uint16(RunStatusRun)
	statsRegs[RegLED] = LEDCelsius | LEDOut1 | LEDAt
	statsRegs[RegPC] = 85  // PC value
	statsRegs[RegPS] = 100 // PS value
	statsRegs[RegPSR] = 95 // 9.5 PSR value

	mock.SetRegisters(0, statsRegs)

	// Configure device info registers
	infoRegs := []uint16{
		0x5058,                 // "PX"
		0x5531,                 // "U1"
		0x3233,                 // "23"
		0x0000, 0x0000, 0x0000, // padding
		125, // firmware 1.25
	}
	mock.SetRegisters(RegInfoStart, infoRegs)

	// Create PXU instance
	pxu, err := NewPxu(6, mock, 2*time.Second, 3)
	if err != nil {
		t.Fatalf("failed to create PXU: %v", err)
	}
	defer pxu.Close()

	// Test reading device stats
	stats, err := pxu.ReadStats()
	if err != nil {
		t.Fatalf("failed to read device stats: %v", err)
	}

	// Validate stats
	if stats.Pv != 23.5 {
		t.Errorf("expected PV 23.5, got %f", stats.Pv)
	}
	if stats.Sp != 25.0 {
		t.Errorf("expected SP 25.0, got %f", stats.Sp)
	}
	if stats.VUnit != "C" {
		t.Errorf("expected unit 'C', got '%s'", stats.VUnit)
	}
	if !stats.Out1 {
		t.Error("expected Out1 to be true")
	}
	if !stats.At {
		t.Error("expected At to be true")
	}
	if stats.RS != RunStatusRun {
		t.Errorf("expected run status RUN, got %v", stats.RS)
	}

	// Test reading device info
	info, err := pxu.ReadInfo()
	if err != nil {
		t.Fatalf("failed to read device info: %v", err)
	}

	// Validate info
	if info.Firmware != "1.25" {
		t.Errorf("expected firmware '1.25', got '%s'", info.Firmware)
	}
}

// TestIntegration_ErrorHandling tests error scenarios
func TestIntegration_ErrorHandling(t *testing.T) {
	mock := NewMockModbus()

	pxu, err := NewPxu(1, mock, 100*time.Millisecond, 2)
	if err != nil {
		t.Fatalf("failed to create PXU: %v", err)
	}

	// Test with connection errors
	mock.SimulateError(true, "connection lost")

	_, err = pxu.ReadStats()
	if err == nil {
		t.Error("expected error when connection is lost")
	}

	_, err = pxu.ReadInfo()
	if err == nil {
		t.Error("expected error when connection is lost")
	}
}

// Benchmark tests
func BenchmarkPxu_ReadStats(b *testing.B) {
	mock := NewMockModbus()

	// Set up realistic register data
	statsRegs := make([]uint16, StatsRegCount)
	statsRegs[RegLED] = LEDCelsius
	mock.SetRegisters(0, statsRegs)

	pxu, err := NewPxu(1, mock, time.Second, 1)
	if err != nil {
		b.Fatalf("failed to create PXU: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := pxu.ReadStats()
		if err != nil {
			b.Fatalf("benchmark failed: %v", err)
		}
	}
}

// Example usage test
func ExamplePxu_ReadStats() {
	// Create a mock for demonstration
	mock := NewMockModbus()

	// Set up some sample data
	registers := make([]uint16, StatsRegCount)
	registers[RegPV] = 235 // 23.5°C
	registers[RegSP] = 250 // 25.0°C
	registers[RegLED] = LEDCelsius | LEDOut1
	registers[RegRS] = uint16(RunStatusRun)
	mock.SetRegisters(0, registers)

	// Create PXU instance
	pxu, err := NewPxu(1, mock, time.Second, 3)
	if err != nil {
		panic(err)
	}
	defer pxu.Close()

	// Read device statistics
	stats, err := pxu.ReadStats()
	if err != nil {
		panic(err)
	}

	fmt.Printf("PV: %.1f%s, SP: %.1f%s, Status: %s\n",
		stats.Pv, stats.VUnit, stats.Sp, stats.VUnit, stats.RS)

	// Output: PV: 23.5C, SP: 25.0C, Status: RUN
}
