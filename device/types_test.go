package device

import (
	"testing"
)

func TestRunStatus_String(t *testing.T) {
	tests := []struct {
		status   RunStatus
		expected string
	}{
		{RunStatusStop, "STOP"},
		{RunStatusRun, "RUN"},
		{RunStatusEnd, "END"},
		{RunStatusPause, "PAUSE"},
		{RunStatusAdvanceProfile, "ADVANCE PROFILE"},
		{RunStatus(99), "UNKNOWN (99)"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.status.String()
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestDeviceInfo_String(t *testing.T) {
	info := Info{
		Model:    "PXU123",
		Firmware: "1.23",
	}

	expected := "Model: PXU123, Firmware: 1.23"
	result := info.String()

	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}
}
