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
	payload := `{"name":"teste","birth_date":"19951001"}`
	tim, err := time.Parse("20060102", "19951001")
	if err != nil {
		t.Error(err)
	}
	v := &Person{Name: "teste", BirthDate: JsonTime(tim)}
	data, err := json.Marshal(v)
	if err != nil {
		t.Error(err)
	}
	if payload != string(data) {
		t.Errorf("expected %v got %v", payload, string(data))
	}
}

func TestJsonTime_UnmarshalJSON(t *testing.T) {
	payload := `{"name": "teste", "birth_date": "19951001"}`
	v := new(Person)
	err := json.Unmarshal([]byte(payload), &v)
	if err != nil {
		t.Error(err)
	}
}
