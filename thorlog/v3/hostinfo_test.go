package thorlog

import (
	"encoding/json"
	"testing"
)

func TestHostInfo_UnmarshalNilPlatform(t *testing.T) {
	for _, data := range []string{`{}`, `{"platform":null}`} {
		if err := json.Unmarshal([]byte(data), &HostInfo{}); err != nil {
			t.Errorf("%s: %v", data, err)
		}
	}
}
