// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package xxml

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// setValueFromCdata sets the given Go value from the given XML cdata.
// The cdata is parsed based on the Go value's type and associated
// rules for the implied XML type.
func setValueFromCdata(fieldValue reflect.Value, val string) (ok bool) {
	ok = false
	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.SetString(val)
		ok = true
	case reflect.Bool:
		fieldValue.SetBool(val == "true" || val == "1")
		ok = true
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		n, _ := strconv.ParseInt(val, 10, 64)
		fieldValue.SetInt(n)
		ok = true
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		n, _ := strconv.ParseUint(val, 10, 64)
		fieldValue.SetUint(n)
		ok = true
	case reflect.Float32, reflect.Float64:
		f, _ := strconv.ParseFloat(val, 64)
		fieldValue.SetFloat(f)
		ok = true
	case reflect.Pointer:
		elemType := fieldValue.Type().Elem()
		newVal := reflect.New(elemType)
		setValueFromCdata(newVal.Elem(), val)
		fieldValue.Set(newVal)
		ok = true
	case reflect.Struct:
		switch fieldValue.Interface().(type) {
		case time.Time:
			time, _ := parseDateTime(val)
			fieldValue.Set(reflect.ValueOf(time))
		}
	}
	return
}

var dateTimeFormats = []string{
	"2006-01-02T15:04:05Z07:00",
	"2006-01-02T15:04:05",
	"2006-01-02Z07:00",
	"2006-01-02",
	"15:04:05Z07:00",
	"15:04:05",
}

func parseDateTime(val string) (time.Time, error) {
	var t time.Time
	var err error
	for _, fmt := range dateTimeFormats {
		t, err = time.Parse(fmt, val)
		if err == nil {
			break
		}
	}
	if err != nil {
		err = fmt.Errorf("cannot parse value as date/time")
	}
	return t, err
}
