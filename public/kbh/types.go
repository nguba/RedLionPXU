package kbh

import (
	"fmt"
	"strings"
)

// MashPlanData represents the top-level structure of your JSON data.
type MashPlanData struct {
	MashPlan []MashStep `json:"Maischplan"`
}

// String provides a human-readable string representation of the MashPlanData.
func (mpd *MashPlanData) String() string {
	var sb strings.Builder
	sb.WriteString("Mash Plan Data:\n")
	for i, step := range mpd.MashPlan {
		sb.WriteString(fmt.Sprintf("--- Step %d ---\n", i+1))
		sb.WriteString(step.String()) // Call the String() method for each MashStep
		if i < len(mpd.MashPlan)-1 {
			sb.WriteString("\n") // Add newline between steps for readability
		}
	}
	return sb.String()
}

// MashStep represents a single step in your mash plan.
// Note: Fields like ExtraDuration1, ExtraDuration2, ExtraTemp1, ExtraTemp2, and WaterTemp
// are defined as interface{} because the provided JSON shows them as both numbers (0) and
// empty strings (""). You will need to perform type assertions when accessing these.
type MashStep struct {
	MashProportion  float64     `json:"AnteilMaische"` // Proportion of mash
	MaltProportion  float64     `json:"AnteilMalz"`    // Proportion of malt
	WaterProportion float64     `json:"AnteilWasser"`  // Proportion of water
	ExtraDuration1  interface{} `json:"DauerExtra1"`   // Extra duration 1 (mixed type: float64 or string)
	ExtraDuration2  interface{} `json:"DauerExtra2"`   // Extra duration 2 (mixed type: float64 or string)
	RestDuration    float64     `json:"DauerRast"`     // Rest duration in minutes
	Name            string      `json:"Name"`          // Name of the mash step
	ExtraTemp1      interface{} `json:"TempExtra1"`    // Extra temperature 1 (mixed type: float64 or string)
	ExtraTemp2      interface{} `json:"TempExtra2"`    // Extra temperature 2 (mixed type: float64 or string)
	MaltTemp        float64     `json:"TempMalz"`      // Malt temperature
	RestTemp        float64     `json:"TempRast"`      // Rest temperature
	WaterTemp       interface{} `json:"TempWasser"`    // Water temperature (mixed type: float64 or string)
	Type            int         `json:"Typ"`           // Type of step
}

// typeNames maps the numerical Type field to its human-readable English name.
// This map can be defined globally or as a field in a struct if its context is specific.
var typeNames = map[int]string{
	0: "Doughing In",
	1: "Heating Up",
	2: "Adding Water", // "Zubruehen" - adding water to a mash or decoction
	3: "Pouring In",   // "Zuschuetten" - likely pouring in something (e.g., more mash, water)
	4: "Decoction",
}

// GetTypeName returns the human-readable English name for the mash step type.
// If the type is unknown, it returns "Unknown Type".
func (ms *MashStep) GetTypeName() string {
	if name, ok := typeNames[ms.Type]; ok {
		return name
	}
	return "Unknown Type"
}

// String provides a human-readable string representation of a single MashStep.
func (ms *MashStep) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("  Name: %s\n", ms.Name))
	sb.WriteString(fmt.Sprintf("  Step Type: %s (Original Type ID: %d)\n", ms.GetTypeName(), ms.Type))
	sb.WriteString(fmt.Sprintf("  Rest Duration: %.0f minutes\n", ms.RestDuration))
	sb.WriteString(fmt.Sprintf("  Rest Temperature: %.2f°C\n", ms.RestTemp))
	sb.WriteString(fmt.Sprintf("  Mash Proportion: %.0f%%\n", ms.MashProportion))
	sb.WriteString(fmt.Sprintf("  Malt Proportion: %.0f%%\n", ms.MaltProportion))
	sb.WriteString(fmt.Sprintf("  Water Proportion: %.2f%%\n", ms.WaterProportion))

	// Handle mixed type fields with type assertions for String() method
	printInterfaceField(&sb, "Extra Duration 1", ms.ExtraDuration1, "%.0f", "Not applicable")
	printInterfaceField(&sb, "Extra Duration 2", ms.ExtraDuration2, "%.0f", "Not applicable")
	printInterfaceField(&sb, "Extra Temperature 1", ms.ExtraTemp1, "%.2f°C", "Not applicable")
	printInterfaceField(&sb, "Extra Temperature 2", ms.ExtraTemp2, "%.2f°C", "Not applicable")
	printInterfaceField(&sb, "Water Temperature", ms.WaterTemp, "%.2f°C", "Not applicable")

	sb.WriteString(fmt.Sprintf("  Malt Temperature: %.2f°C\n", ms.MaltTemp))

	return sb.String()
}

// Helper function to print interface{} fields gracefully
func printInterfaceField(sb *strings.Builder, fieldName string, value interface{}, format, naMessage string) {
	if val, ok := value.(float64); ok {
		sb.WriteString(fmt.Sprintf("  %s: %s\n", fieldName, fmt.Sprintf(format, val)))
	} else if val, ok := value.(string); ok && val == "" {
		sb.WriteString(fmt.Sprintf("  %s: %s\n", fieldName, naMessage))
	} else {
		// Fallback for unexpected types
		sb.WriteString(fmt.Sprintf("  %s (unknown type): %v\n", fieldName, value))
	}
}
