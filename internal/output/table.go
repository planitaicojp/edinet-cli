package output

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"
)

type TableFormatter struct{}

func (f *TableFormatter) Format(w io.Writer, data any) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	defer func() { _ = tw.Flush() }()

	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Slice {
		if v.Len() == 0 {
			_, _ = fmt.Fprintln(w, "No results.")
			return nil
		}
		elemType := v.Type().Elem()
		headers := structHeaders(elemType)
		_, _ = fmt.Fprintln(tw, strings.Join(headers, "\t"))

		for i := 0; i < v.Len(); i++ {
			vals := structValues(v.Index(i))
			_, _ = fmt.Fprintln(tw, strings.Join(vals, "\t"))
		}
		return nil
	}

	if v.Kind() == reflect.Struct {
		headers := structHeaders(v.Type())
		_, _ = fmt.Fprintln(tw, strings.Join(headers, "\t"))
		vals := structValues(v)
		_, _ = fmt.Fprintln(tw, strings.Join(vals, "\t"))
		return nil
	}

	_, err := fmt.Fprintln(w, data)
	return err
}

func structHeaders(t reflect.Type) []string {
	var headers []string
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		name := strings.Split(tag, ",")[0]
		headers = append(headers, strings.ToUpper(name))
	}
	return headers
}

func structValues(v reflect.Value) []string {
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
