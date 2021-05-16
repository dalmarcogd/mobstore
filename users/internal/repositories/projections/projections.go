package projections

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func GetProjections(st interface{}) []string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error getting projections err=%v", r)
		}
	}()

	t := reflect.TypeOf(st)
	v := reflect.ValueOf(st)

	projections := make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.FieldByName(field.Name)
		// Get the field tag value
		if tag, ok := field.Tag.Lookup("projection"); ok && !fieldValue.IsZero() && fieldValue.CanInterface() {
			if fieldValue.Interface() == true {
				projections = append(projections, tag)
			}
		}
	}
	return projections
}

func GetProjectionValue(st interface{}, projection string) interface{} {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error getting projection err=%v", r)
		}
	}()

	t := reflect.TypeOf(st)
	v := reflect.ValueOf(st)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.FieldByName(field.Name)
		// Get the field tag value
		if tag, ok := field.Tag.Lookup("projection"); ok && tag == projection {
			if fieldValue.Type().Kind() == reflect.Ptr {
				return fieldValue.Elem().Interface()
			} else {
				return fieldValue.Interface()
			}
		}
	}
	return nil
}

func SetProjections(st interface{}, val map[string]interface{}) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error getting projections err=%v", r)
		}
	}()
	valMap := make(map[string]interface{})
	for k, v := range val {
		valMap[k] = v
	}

	v := reflect.ValueOf(st)
	if v.Kind() != reflect.Ptr {
		return errors.New("st should be pointer")
	}
	v = v.Elem()
	t := reflect.TypeOf(st).Elem()

	for column, value := range valMap {
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			fieldValue := v.FieldByName(field.Name)
			if tag, ok := field.Tag.Lookup("projection"); ok && tag == column && fieldValue.CanInterface() {
				sqlVal := reflect.ValueOf(value).Interface().(sql.RawBytes)
				switch fieldValue.Kind() {
				case reflect.String:
					fieldValue.SetString(string(sqlVal))
				case reflect.Ptr:
					fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
					switch fieldValue.Type().Elem().Kind() {
					case reflect.String:
						fieldValue.Elem().SetString(string(sqlVal))
					case reflect.Int64:
						parseInt, err := strconv.ParseInt(string(sqlVal), 10, 64)
						if err != nil {
							return err
						}
						fieldValue.Elem().SetInt(parseInt)
					default:
						fieldValue.Elem().Set(reflect.ValueOf(value))
					}
				default:
					fieldValue.Set(reflect.ValueOf(value))

				}
				delete(valMap, column)
			}
		}
	}
	if len(valMap) > 0 {
		return errors.New("there are unassignable columns")
	}
	return nil
}
