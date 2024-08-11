package customtypes

import (
	"github.com/goccy/go-reflect"
	"time"
)

type Time time.Time

var TimeLayout = "2006-01-02 15:04:05"

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	if data == nil || len(data) == 0 {
		*t = Time(time.Time{})
		return nil
	}
	if data[0] == '"' && data[len(data)-1] == '"' {
		data = data[1 : len(data)-1]
	}
	if data == nil || len(data) == 0 {
		*t = Time(time.Time{})
		return nil
	}
	if string(data) == "null" {
		*t = Time(time.Time{})
		return nil
	}
	now, err := time.ParseInLocation(TimeLayout, string(data), time.Local)
	*t = Time(now)
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	if reflect.ValueOf(t).IsZero() {
		return []byte("null"), nil
	}
	b := make([]byte, 0, len(TimeLayout)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TimeLayout)
	b = append(b, '"')
	return b, nil
}

func (t Time) String() string {
	return time.Time(t).Format(TimeLayout)
}
