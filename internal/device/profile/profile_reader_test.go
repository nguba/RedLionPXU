package profile

import (
	"fmt"
	"io"
	"log"
	"os"
	"testing"
)

func TestJSONProfileReader_ReadProfiles(t *testing.T) {
	type fields struct {
		dataSource io.Reader
	}
	//tests := []struct {
	//	name    string
	//	fields  fields
	//	want    []Profile
	//	wantErr bool
	//}{
	//	// TODO: Add test cases.
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		jpr := &JSONProfileReader{
	//			dataSource: tt.fields.dataSource,
	//		}
	//		got, err := jpr.ReadProfiles()
	//		if (err != nil) != tt.wantErr {
	//			t.Errorf("ReadProfiles() error = %v, wantErr %v", err, tt.wantErr)
	//			return
	//		}
	//		if !reflect.DeepEqual(got, tt.want) {
	//			t.Errorf("ReadProfiles() got = %v, want %v", got, tt.want)
	//		}
	//	})
	//}

	file, err := os.Open("profiles.json")
	if err != nil {
		log.Fatal(err)
	}
	reader := NewJSONProfileReader(file)

	// Read the profiles using the interface method
	profiles, err := reader.ReadProfiles()
	if err != nil {
		fmt.Printf("Error reading profiles: %v\n", err)
		return
	}

	// Print the parsed profiles to verify
	fmt.Println("Successfully parsed brewing profiles:")
	for _, p := range profiles {
		fmt.Printf("  Profile ID: %d\n  Profile Name: %s\n  Link to Next Profile: %d\n  Segments:\n", p.ID, p.ProfileName, p.LinkToNextProfile)
		for _, s := range p.Segments {
			fmt.Printf("    - Segment ID: %d, Setpoint: %.1f°C, Time: %.1f minutes, Notes: %s\n", s.ID, s.Setpoint, s.Time, s.Notes)
		}
		fmt.Println() // Add a newline for readability between profiles
	}

	// Example of how you might interact with the data by ID
	profileMap := make(map[int]Profile)
	for _, p := range profiles {
		profileMap[p.ID] = p
	}

	if p0, ok := profileMap[0]; ok {
		fmt.Printf("Example: Accessed Profile 0 by ID. Name: %s\n", p0.ProfileName)
	}
	if p1, ok := profileMap[1]; ok {
		fmt.Printf("Example: Accessed Profile 1 by ID. Name: %s\n", p1.ProfileName)
	}
}
