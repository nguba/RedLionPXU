package kbh

import (
	"encoding/json"
	"fmt"
	"io"
)

type KbhReader struct {
	dataSource io.Reader
}

func NewKbhReader(dataSource io.Reader) *KbhReader {
	return &KbhReader{dataSource}
}

func (r *KbhReader) ReadRecipe() (*MashPlanData, error) {
	var data MashPlanData
	decoder := json.NewDecoder(r.dataSource)

	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("error decoding recipe: %v", err)
	}

	return &data, nil
}
