package kbh

import (
	"os"
	"testing"
)

func TestKbhReader_ReadRecipe(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		want     *MashPlanData
		wantErr  bool
	}{
		{"load recipe from JSON", "sample_export.json", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			file, err := os.Open(tt.fileName)
			if err != nil {
				t.Fatal(err)
			}
			r := &KbhReader{
				dataSource: file,
			}
			_, err = r.ReadRecipe()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadRecipe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("ReadRecipe() got = %v, want %v", got, tt.want)
			//}
		})
	}
}
