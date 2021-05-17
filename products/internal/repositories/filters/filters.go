package filters

import (
	"fmt"
	"reflect"
)

func GetFilters(ts interface{}) map[string]interface{} {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error getting filters err=%v", r)
		}
	}()

	t := reflect.TypeOf(ts)
	v := reflect.ValueOf(ts)

	filters := make(map[string]interface{}, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.FieldByName(field.Name)
		// Get the field tag value
		if tag, ok := field.Tag.Lookup("filter"); ok && !fieldValue.IsZero() && fieldValue.CanInterface() {
			if _, ok := filters[tag]; !ok {
				if fieldValue.Type().Kind() == reflect.Ptr {
					filters[tag] = fieldValue.Elem().Interface()
				} else {
					filters[tag] = fieldValue.Interface()
				}
			} else {
				panic("there filter tag duplicate")
			}
		}
	}
	return filters
}
