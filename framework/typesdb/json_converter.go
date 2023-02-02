package typesdb

import "time"

type Datetime struct {
	time.Time
}

func (m *Datetime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}

	tt, err := time.Parse(`"`+time.RFC3339+`"`, string(data))
	*m = Datetime{tt}
	return err
}
