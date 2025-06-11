package public

// Segment represents a single step within a brewing profile.
type Segment struct {
	ID       int     `json:"id"`
	Setpoint float64 `json:"setpoint"`
	Time     float64 `json:"time"` // Duration of the segment in minutes (max 999.9)
	Notes    string  `json:"notes"`
}

// Profile represents a complete brewing phase with multiple segments.
type Profile struct {
	ID                int       `json:"id"`
	ProfileName       string    `json:"profile_name"`
	LinkToNextProfile int       `json:"link_to_next_profile"` // Integer ID of the next profile, or -1 to indicate end
	Segments          []Segment `json:"segments"`
}

// ProfileData is the top-level structure holding all brewing profiles.
type ProfileData struct {
	Profiles []Profile `json:"profiles"`
}
