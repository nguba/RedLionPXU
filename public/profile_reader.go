package public

import (
	"encoding/json"
	"fmt"
	"io"
)

// JSONProfileReader implements ProfileReader for JSON data.
type JSONProfileReader struct {
	dataSource io.Reader
}

// NewJSONProfileReader creates a new JSONProfileReader instance.
// It takes an io.Reader, which can be a file, a network connection, or a bytes.Buffer.
func NewJSONProfileReader(r io.Reader) *JSONProfileReader {
	return &JSONProfileReader{dataSource: r}
}

// ReadProfiles reads and parses brewing profiles from the underlying data source.
func (jpr *JSONProfileReader) ReadProfiles() ([]Profile, error) {
	var data ProfileData
	decoder := json.NewDecoder(jpr.dataSource)

	// Decode the JSON data into the ProfileData struct
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode profiles JSON: %w", err)
	}

	return data.Profiles, nil
}
