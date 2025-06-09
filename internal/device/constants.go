package device

import "time"

// General Register addresses
const (
	RegPV     = 0  // Process Value
	RegSP     = 1  // Active Setpoint
	RegTP     = 10 // Proportional Band
	RegTI     = 11 // Integral Time
	RegTD     = 12 // Derivative Time
	RegTGroup = 14 // Parameter Set Selection
	RegRS     = 17 // Controller Status
	RegLED    = 20 // LED Status
	RegPC     = 25 // Current Profile
	RegPS     = 26 // Current Profile Segment
	RegPSR    = 27 // Profile Segment Remaining Time

	StatsRegCount = 30
)

// Device info registers
const (
	RegInfoStart = 1000

	InfoRegCount = 7
)

// Profile registers
const (
	RegProfDEV              = 1090
	RegProfEBT              = 1091
	RegProfIRR              = 1092
	RegProfSegmentStart     = 1100
	RegProfSegmentCount     = 32 // setpoint -> even idx, time -> odd idx
	RegNumSegmentsStart     = 1630
	RegNumSegmentsCount     = 15
	RegProfCycleRepeatStart = 1650
	RegProfCycleRepeatCount = 15
	RegProfLink             = 1670
	RegProfLinkCount        = 15
)

const ()

// LED status bit masks
const (
	LEDAt         uint16 = 1 << 7 // Auto-Tune On
	LEDOut1       uint16 = 1 << 6 // Output Power1 Active
	LEDOut2       uint16 = 1 << 5 // Output Power2 Active
	LEDCelsius    uint16 = 1 << 3 // Display Units are in Celsius
	LEDFahrenheit uint16 = 1 << 2 // Display Units are in Fahrenheit
)

// Default configuration values
const (
	DefaultTimeout = 5 * time.Second
	DefaultRetries = 3
	DefaultSpeed   = 38400
)
