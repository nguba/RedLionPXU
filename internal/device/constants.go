package device

import "time"

// Register addresses
const (
	StatsRegPV     = 0
	StatsRegSP     = 1
	StatsRegTP     = 10
	StatsRegTI     = 11
	StatsRegTD     = 12
	StatsRegTGroup = 14
	StatsRegRS     = 17
	StatsRegLED    = 20
	StatsRegPC     = 25
	StatsRegPS     = 26
	StatsRegPSR    = 27

	StatsRegCount = 30
)

// Device info registers
const (
	InfoRegStart = 1000

	InfoRegCount = 7
)

// Profile registers
const (
	ProfDeviationErrVal     = 1090
	ProfErrorBandTimeout    = 1091
	ProfInitialRampRate     = 1092
	ProfSegmentRegStart     = 1100
	ProfSegmentCount        = 32 // setpoint -> odd num, time -> even num
	ProfNumSegmentsRegStart = 1630
	ProfNumSegmentsCount    = 15
	ProfCycleRepeatStart    = 1650
	ProfCycleRepeatCount    = 15
	ProfLinkProfile         = 1670
	ProfLinkProfileCount    = 15
)

const ()

// LED status bit masks
const (
	LEDAt         uint16 = 1 << 7
	LEDOut1       uint16 = 1 << 6
	LEDOut2       uint16 = 1 << 5
	LEDCelsius    uint16 = 1 << 3
	LEDFahrenheit uint16 = 1 << 2
)

// Default configuration values
const (
	DefaultTimeout = 5 * time.Second
	DefaultRetries = 3
	DefaultSpeed   = 38400
)
