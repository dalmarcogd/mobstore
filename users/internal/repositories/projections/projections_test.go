package projections

import (
	"database/sql"
	"reflect"
	"testing"
)

func TestGetProjections(t *testing.T) {
	type args struct {
		ts interface{}
	}

	type test struct {
		name string
		args args
		want []string
	}
	type testProjection struct {
		Name   bool `projection:"Name"`
		Value1 bool `projection:"Value"`
		Value2 bool `projection:"Amount"`
	}

	var tests []test
	tests = append(tests, test{name: "GetProjections-1", args: args{ts: testProjection{Name: true, Value1: true, Value2: false}}, want: []string{"Name", "Value"}})
	tests = append(tests, test{name: "GetProjections-2", args: args{ts: testProjection{Name: false, Value1: true, Value2: false}}, want: []string{"Value"}})
	tests = append(tests, test{name: "GetProjections-3", args: args{ts: testProjection{}}, want: []string{}})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetProjections(tt.args.ts); got != nil {
				got2 := make(map[string]bool, 0)
				for _, s := range got {
					got2[s] = true
				}
				for _, vw := range tt.want {
					for k := range got2 {
						if k == vw {
							delete(got2, k)
							break
						}
					}
				}
				if len(got2) > 0 {
					t.Errorf("GetProjections() unexpected value = %v, want %v", got2, tt.want)
				}
			}
		})
	}
}

func TestSetProjections(t *testing.T) {
	type args struct {
		st  interface{}
		val map[string]interface{}
	}
	type testProjection struct {
		Name  string `projection:"Name"`
		Name1 string `projection:"Name"`
		Name2 string `projection:"Name2"`
	}
	type test struct {
		name    string
		args    args
		wantErr bool
	}
	var tests []test
	projection := &testProjection{}
	tests = append(tests, test{name: "SetProjections-1", args: args{st: projection, val: map[string]interface{}{"Name": sql.RawBytes("123")}}, wantErr: false})
	tests = append(tests, test{name: "SetProjections-2", args: args{st: projection, val: map[string]interface{}{"Value": sql.RawBytes("33")}}, wantErr: true})
	tests = append(tests, test{name: "SetProjections-3", args: args{st: projection, val: map[string]interface{}{}}, wantErr: false})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetProjections(tt.args.st, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("SetProjections() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				valueOf := reflect.ValueOf(tt.args.st).Elem()
				typeOf := reflect.TypeOf(tt.args.st).Elem()
				for column, value := range tt.args.val {
					for i := 0; i < valueOf.NumField(); i++ {
						field := typeOf.Field(i)
						fieldValue := valueOf.FieldByName(field.Name)
						if tag, ok := field.Tag.Lookup("projection"); ok && tag == column && fieldValue.CanInterface() && fieldValue.Interface() != string(value.(sql.RawBytes)) {
							t.Errorf("SetProjections() expected equal values from val st= %v, val %v", tt.args.st, tt.args.val)
						}
					}
				}
			}
		})
	}
}
