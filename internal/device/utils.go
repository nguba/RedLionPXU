package device

import (
	"encoding/hex"
	"fmt"
	"math"
)

func makeProfile(regs []uint16) string {
	s1 := toFloat(regs[0])
	s2 := toFloat(regs[1])
	s3 := toFloat(regs[2])
	s4 := toFloat(regs[3])
	s5 := toFloat(regs[4])
	s6 := toFloat(regs[5])
	return fmt.Sprintf("%.2f - %.2f\n%.2f - %.2f\n%.2f - %.2f\n", s1, s2, s3, s4, s5, s6)
}

func toFloat(r uint16) float64 {
	return float64(r) / 10
}

func toUint16(r float64) uint16 {
	scaled := r * 10.0
	truncated := math.Floor(scaled)

	return uint16(truncated)
}

func toString(input uint16) string {
	r := fmt.Sprintf("%04x", input) // Ensure 4 digits with leading zeros
	bs, err := hex.DecodeString(r)
	if err != nil {
		return ""
	}

	// Filter out null bytes and control characters
	result := make([]byte, 0, len(bs))
	for _, b := range bs {
		if b >= 32 && b <= 126 { // Printable ASCII range
			result = append(result, b)
		}
	}

	return string(result)
}
