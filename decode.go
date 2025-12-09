package cnab

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func decode(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrInvalidStructPtr
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return ErrInvalidStruct
	}

	line := string(data)

	currentPos := 0
	t := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		field := t.Field(i)
		tagValue := field.Tag.Get("cnab")
		if tagValue == "" {
			continue
		}

		tag, err := parseTag(tagValue)
		if err != nil {
			return fmt.Errorf("field %s: %w", field.Name, err)
		}

		// Determine start and end positions
		start := currentPos
		if tag.start > 0 {
			start = tag.start - 1 // 1-based to 0-based
		}

		end := start + tag.size

		if end > len(line) {
			return ErrLineTooShort
		}

		valStr := line[start:end]
		// Update currentPos for next field if using sequential
		currentPos = end

		err = setFieldValue(rv.Field(i), valStr, tag)
		if err != nil {
			return fmt.Errorf("field %s: %w", field.Name, err)
		}
	}

	return nil
}

func setFieldValue(v reflect.Value, s string, tag fieldTag) error {
	// Check for Unmarshaler interface
	// v is likely addressable since it comes from rv.Field(i) of a pointer struct
	if v.CanAddr() {
		addr := v.Addr()
		if u, ok := addr.Interface().(Unmarshaler); ok {
			return u.UnmarshalCNAB([]byte(s))
		}
	}

	// Also check if the value itself implements Unmarshaler (unlikely for pointer receivers, but possible)
	if v.CanInterface() {
		if u, ok := v.Interface().(Unmarshaler); ok {
			return u.UnmarshalCNAB([]byte(s))
		}
	}

	switch v.Kind() {
	case reflect.String:
		if tag.align == "right" {
			s = strings.TrimLeft(s, string(tag.fill))
		} else {
			s = strings.TrimRight(s, string(tag.fill))
		}
		v.SetString(s)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s = strings.TrimSpace(s)
		if s == "" {
			v.SetInt(0)
			return nil
		}
		val, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s = strings.TrimSpace(s)
		if s == "" {
			v.SetUint(0)
			return nil
		}
		val, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(val)
	case reflect.Float32, reflect.Float64:
		s = strings.TrimSpace(s)
		if s == "" {
			v.SetFloat(0)
			return nil
		}

		if tag.decimal > 0 {
			// Parse as int first, then divide
			valInt, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				// Fallback to float parsing if it has a dot?
				f, err2 := strconv.ParseFloat(s, 64)
				if err2 == nil {
					v.SetFloat(f)
					return nil
				}
				return err
			}
			f := float64(valInt) / math.Pow(10, float64(tag.decimal))
			v.SetFloat(f)
		} else {
			val, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return err
			}
			v.SetFloat(val)
		}

	case reflect.Struct:
		if v.Type() == reflect.TypeOf(time.Time{}) {
			s = strings.TrimSpace(s)
			if s == "" || s == strings.Repeat(string(tag.fill), len(s)) {
				// Empty time
				return nil
			}
			f := tag.format
			if f == "" {
				f = "20060102"
			}
			t, err := time.Parse(f, s)
			if err != nil {
				return ErrInvalidDateFormat
			}
			v.Set(reflect.ValueOf(t))
		}
	default:
		return ErrUnsupportedType
	}
	return nil
}
