package cnab

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func encode(v interface{}) ([]byte, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return nil, ErrInvalidStruct
	}

	var buf bytes.Buffer
	t := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		field := t.Field(i)
		tagValue := field.Tag.Get("cnab")
		if tagValue == "" {
			continue
		}

		tag, err := parseTag(tagValue)
		if err != nil {
			return nil, fmt.Errorf("field %s: %w", field.Name, err)
		}

		val := rv.Field(i)
		s, err := formatValue(val, tag)
		if err != nil {
			return nil, fmt.Errorf("field %s: %w", field.Name, err)
		}

		if len(s) > tag.size {
			return nil, fmt.Errorf("field %s: value '%s' too long for size %d", field.Name, s, tag.size)
		}

		// Apply padding
		padding := tag.size - len(s)
		if padding > 0 {
			padStr := strings.Repeat(string(tag.fill), padding)
			if tag.align == "right" {
				if tag.fill == '0' && strings.HasPrefix(s, "-") {
					// Handle negative number with zero padding: -0001
					s = "-" + padStr + s[1:]
				} else {
					s = padStr + s
				}
			} else {
				s = s + padStr
			}
		}

		buf.WriteString(s)
	}

	return buf.Bytes(), nil
}

func formatValue(v reflect.Value, tag fieldTag) (string, error) {
	// Check for Marshaler interface
	if v.CanInterface() {
		if m, ok := v.Interface().(Marshaler); ok {
			b, err := m.MarshalCNAB()
			if err != nil {
				return "", err
			}
			return string(b), nil
		}
	}
	
	// Also check if addressable and pointer implements Marshaler
	// (e.g. func (t *Type) MarshalCNAB...)
	if v.CanAddr() {
		if m, ok := v.Addr().Interface().(Marshaler); ok {
			b, err := m.MarshalCNAB()
			if err != nil {
				return "", err
			}
			return string(b), nil
		}
	}

	switch v.Kind() {
	case reflect.String:
		return v.String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10), nil
	case reflect.Float32, reflect.Float64:
		if tag.decimal > 0 {
			// Multiply by 10^decimal and round to int
			f := v.Float()
			mult := math.Pow(10, float64(tag.decimal))
			val := int64(math.Round(f * mult))
			return strconv.FormatInt(val, 10), nil
		}
		return strconv.FormatFloat(v.Float(), 'f', -1, 64), nil
	case reflect.Struct:
		if v.Type() == reflect.TypeOf(time.Time{}) {
			t := v.Interface().(time.Time)
			if t.IsZero() {
				return "", nil
			}
			f := tag.format
			if f == "" {
				f = "20060102" // Default CNAB date format
			}
			return t.Format(f), nil
		}
	}
	return "", ErrUnsupportedType
}
