package customtypes

import "time"

type Time time.Time

var TimeLayout = "2006-01-02 15:04:05"

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+TimeLayout+`"`, string(data), time.Local)
	*t = Time(now)
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TimeLayout)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TimeLayout)
	b = append(b, '"')
	return b, nil
}

func (t Time) String() string {
	return time.Time(t).Format(TimeLayout)
}
