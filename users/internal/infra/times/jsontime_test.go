package times

import (
	"encoding/json"
	"testing"
	"time"
)

type Person struct {
	Name      string   `json:"name"`
	BirthDate JsonTime `json:"birth_date"`
}

func TestJsonTime_MarshalJSON(t *testing.T) {
	payload := `{"name": "teste", "birth_date": "1995-10-01"}`
	tim, err := time.Parse("2006-01-02", "1995-10-01")
	if err != nil {
		t.Error(err)
	}
	v := &Person{Name: "teste", BirthDate: JsonTime(tim)}
	data, err := json.Marshal(v)
	if err != nil {
		t.Error()
	}
	if payload != string(data) {
		t.Errorf("expected %v got %v", payload, data)
	}
}

func TestJsonTime_UnmarshalJSON(t *testing.T) {
	payload := `{"name": "teste", "birth_date": "1995-10-01"}`
	v := new(Person)
	err := json.Unmarshal([]byte(payload), &v)
	if err != nil {
		t.Error()
	}
}
