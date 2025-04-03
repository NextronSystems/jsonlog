package thorlog

import (
	"encoding/json"
	"time"
)

type Time struct {
	time.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time.UnixMilli())
}

func (t *Time) UnmarshalJSON(data []byte) error {
	var unixMilli int64
	if err := json.Unmarshal(data, &unixMilli); err != nil {
		return err
	}
	t.Time = time.UnixMilli(unixMilli)
	return nil
}
