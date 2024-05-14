package common

import (
	"encoding/json"
	"testing"
)

func TestVersion_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		v       Version
		args    []byte
		wantErr bool
	}{
		{
			name:    "String version",
			v:       "v1.0.1",
			args:    []byte(`"v1.0.1"`),
			wantErr: false,
		},
		{
			name:    "Invalid string version",
			v:       "",
			args:    []byte(`"not a version"`),
			wantErr: true,
		},
		{
			name:    "Int version",
			v:       "v1.0.0",
			args:    []byte(`1`),
			wantErr: false,
		},
		{
			name:    "Invalid version type",
			v:       "",
			args:    []byte(`[]`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var version Version
			err := json.Unmarshal(tt.args, &version)
			if (err != nil) != tt.wantErr {
				t.Errorf("Version.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && version != tt.v {
				t.Errorf("Version.UnmarshalJSON() got = %v, want %v", version, tt.v)
			}
		})
	}
}
