package events

import (
	"encoding/json"
	"time"
)

const timeFormat = "2006-01-02T15:04:05.999"

type Time struct {
	time.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Format(timeFormat))
}

func (t *Time) UnmarshalJSON(data []byte) error {
	var timeStr string
	if err := json.Unmarshal(data, &timeStr); err != nil {
		return err
	}
	parsed, err := time.Parse(timeFormat, timeStr)
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}
