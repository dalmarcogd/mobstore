package times

import (
	"fmt"
	"strings"
	"time"
)

type JsonTime time.Time

func (j *JsonTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("20060102", s)
	if err != nil {
		return err
	}
	*j = JsonTime(t)
	return nil
}

func (j JsonTime) MarshalJSON() ([]byte, error) {
	if t := time.Time(j); !t.IsZero() {
		return []byte(fmt.Sprintf("\"%v\"", t.Format("20060102"))), nil
	}
	return nil, nil
}
