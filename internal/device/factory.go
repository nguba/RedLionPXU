package device

import (
	"fmt"
	"strings"
	"time"
)

// DefaultConfiguration returns a default configuration for COM3
func DefaultConfiguration() *Configuration {
	return &Configuration{
		URL:      "rtu://COM3",
		Speed:    DefaultSpeed,
		DataBits: 8,
		Parity:   "none",
		Timeout:  500 * time.Millisecond,
	}
}

// NewDefaultPxu creates a PXU with default settings
func NewDefaultPxu(unitId uint8) (*Pxu, error) {
	cfg := DefaultConfiguration()
	client, err := NewModbusHandler(cfg)
	if err != nil {
		return nil, err
	}

	return NewPxu(unitId, client, DefaultTimeout, DefaultRetries)
}

func NewStats(regs []uint16) (*Stats, error) {
	ledStatus := regs[RegLED]

	unit, err := parseTemperatureUnit(ledStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to parse temperature unit: %w", err)
	}

	return &Stats{
		Pv:     toFloat(regs[RegPV]),
		Sp:     toFloat(regs[RegSP]),
		Out1:   ledStatus&LEDOut1 != 0,
		Out2:   ledStatus&LEDOut2 != 0,
		At:     ledStatus&LEDAt != 0,
		TP:     toFloat(regs[RegTP]),
		TI:     regs[RegTI],
		TD:     regs[RegTD],
		TGroup: regs[RegTGroup],
		RS:     RunStatus(regs[RegRS]),
		VUnit:  unit,
		PC:     regs[RegPC],
		PS:     regs[RegPS],
		PSR:    toFloat(regs[RegPSR]),
	}, nil
}

func parseTemperatureUnit(ledStatus uint16) (string, error) {
	celsius := ledStatus&LEDCelsius != 0
	fahrenheit := ledStatus&LEDFahrenheit != 0

	switch {
	case celsius && fahrenheit:
		return "", fmt.Errorf("invalid temperature unit: both flags set (0x%04X)", ledStatus)
	case celsius:
		return "C", nil
	case fahrenheit:
		return "F", nil
	default:
		return "", fmt.Errorf("no temperature unit specified in LED status: 0x%04X", ledStatus)
	}
}

func NewInfo(regs []uint16) (*Info, error) {
	var model strings.Builder
	l := InfoRegCount - 1
	for i := 0; i < l; i++ {
		model.WriteString(toString(regs[i]))
	}

	firmware := fmt.Sprintf("%.2f", float64(regs[l])/100)

	return &Info{
		Model:    strings.TrimSpace(model.String()),
		Firmware: firmware,
	}, nil
}
