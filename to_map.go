package json

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// ToMap converts any struct to map[string]any, preserving numeric types
// (int64 stays int64, float64 stays float64, no precision loss).
// Handles nested structs, pointers, slices, maps, and any/interface{} fields recursively.
func ToMap(s any) map[string]any {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	return structToMap(v)
}

func structToMap(v reflect.Value) map[string]any {
	t := v.Type()
	m := make(map[string]any, t.NumField())

	for i := range t.NumField() {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}

		name := fieldName(f)
		m[name] = toMapValue(v.Field(i))
	}
	return m
}

func fieldName(f reflect.StructField) string {
	tag := f.Tag.Get("json")
	if tag == "" || tag == "-" {
		return f.Name
	}
	if name, _, _ := strings.Cut(tag, ","); name != "" {
		return name
	}
	return tag
}

func toMapValue(v reflect.Value) any {
	switch v.Kind() {
	case reflect.Pointer:
		if v.IsNil() {
			return nil
		}
		return toMapValue(v.Elem())

	case reflect.Interface:
		if v.IsNil() {
			return nil
		}
		return toMapValue(v.Elem())

	case reflect.Struct:
		// time.Time is a special case: marshal as string, not a map
		if _, ok := v.Interface().(time.Time); ok {
			return v.Interface()
		}
		return structToMap(v)

	case reflect.Slice, reflect.Array:
		n := v.Len()
		s := make([]any, n)
		for i := range n {
			s[i] = toMapValue(v.Index(i))
		}
		return s

	case reflect.Map:
		m := make(map[string]any, v.Len())
		for _, key := range v.MapKeys() {
			m[anyToString(key)] = toMapValue(v.MapIndex(key))
		}
		return m
	default:
		return v.Interface()
	}
}

func anyToString(v reflect.Value) string {
	if v.Kind() == reflect.String {
		return v.String()
	}
	return fmt.Sprintf("%v", v.Interface())
}
