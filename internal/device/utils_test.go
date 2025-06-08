package device

import (
	"testing"
)

func TestToFloat(t *testing.T) {
	tests := []struct {
		input    uint16
		expected float64
	}{
		{0, 0.0},
		{10, 1.0},
		{123, 12.3},
		{65535, 6553.5},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := toFloat(tt.input)
			if result != tt.expected {
				t.Errorf("toFloat(%d) = %f, expected %f", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		name     string
		input    uint16
		expected string
	}{
		{"ASCII letters", 0x4142, "AB"},
		{"ASCII numbers", 0x3132, "12"},
		{"Mixed", 0x4131, "A1"},
		{"With nulls", 0x4100, "A"},    // Should filter out null bytes
		{"Control chars", 0x0141, "A"}, // Should filter out control chars
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toString(tt.input)
			if result != tt.expected {
				t.Errorf("toString(0x%04X) = '%s', expected '%s'", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMakeProfile(t *testing.T) {
	t.SkipNow()
	tests := []struct {
		name     string
		input    []uint16
		expected string
	}{
		{
			name:     "valid profile",
			input:    []uint16{100, 200, 300, 400, 500, 600},
			expected: "10.00 - 20.00\n30.00 - 40.00\n50.00 - 60.00\n",
		},
		{
			name:     "insufficient data",
			input:    []uint16{100, 200},
			expected: "",
		},
		{
			name:     "empty input",
			input:    []uint16{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := makeProfile(tt.input)
			if result != tt.expected {
				t.Errorf("makeProfile() = '%s', expected '%s'", result, tt.expected)
			}
		})
	}
}
