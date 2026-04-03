package output

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strings"
)

type CSVFormatter struct{}

func (f *CSVFormatter) Format(w io.Writer, data any) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Slice {
		if v.Len() == 0 {
			return nil
		}
		elemType := v.Type().Elem()
		headers := csvHeaders(elemType)
		if err := cw.Write(headers); err != nil {
			return err
		}
		for i := 0; i < v.Len(); i++ {
			vals := csvValues(v.Index(i))
			if err := cw.Write(vals); err != nil {
				return err
			}
		}
		return nil
	}

	if v.Kind() == reflect.Struct {
		headers := csvHeaders(v.Type())
		if err := cw.Write(headers); err != nil {
			return err
		}
		return cw.Write(csvValues(v))
	}

	return cw.Write([]string{fmt.Sprint(data)})
}

func csvHeaders(t reflect.Type) []string {
	var headers []string
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		headers = append(headers, strings.Split(tag, ",")[0])
	}
	return headers
}

func csvValues(v reflect.Value) []string {
	var vals []string
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		vals = append(vals, fmt.Sprintf("%v", v.Field(i).Interface()))
	}
	return vals
}
