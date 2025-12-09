package dynamic

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// Marshal takes a map of data and a list of fields definition, returning a CNAB line.
func Marshal(data map[string]interface{}, layout []Field) ([]byte, error) {
	var buf bytes.Buffer

	for _, field := range layout {
		val, ok := data[field.Name]
		if !ok && field.Required {
			return nil, fmt.Errorf("field %s is required but missing", field.Name)
		}

		s, err := formatValue(val, field)
		if err != nil {
			return nil, fmt.Errorf("field %s: %w", field.Name, err)
		}

		if len(s) > field.Size {
			return nil, fmt.Errorf("field %s: value '%s' too long for size %d", field.Name, s, field.Size)
		}

		// Defaults
		fill := " "
		if field.Fill != "" {
			fill = field.Fill
		} else if field.Type == "int" || field.Type == "float" {
			fill = "0"
		}

		align := "left"
		if field.Align != "" {
			align = field.Align
		} else if field.Type == "int" || field.Type == "float" {
			align = "right"
		}

		// Apply padding
		padding := field.Size - len(s)
		if padding > 0 {
			padStr := strings.Repeat(fill, padding)
			if align == "right" {
				if fill == "0" && strings.HasPrefix(s, "-") {
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

func formatValue(v interface{}, f Field) (string, error) {
	if v == nil {
		return "", nil
	}

	switch val := v.(type) {
	case string:
		return formatStringValue(val, f)
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", val), nil
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", val), nil
	case float32, float64:
		vf := 0.0
		if f64, ok := val.(float64); ok {
			vf = f64
		} else {
			vf = float64(val.(float32))
		}

		if f.Decimal > 0 {
			mult := math.Pow(10, float64(f.Decimal))
			vi := int64(math.Round(vf * mult))
			return fmt.Sprintf("%d", vi), nil
		}
		// Default float format if decimal not specified? Or error?
		// Let's assume standard float string
		return strconv.FormatFloat(vf, 'f', -1, 64), nil
	case time.Time:
		format := f.Format
		if format == "" {
			format = "20060102"
		}
		return val.Format(format), nil
	}

	return fmt.Sprintf("%v", v), nil
}

// formatStringValue coerces string values according to the declared field type,
// returning descriptive errors when conversion is not possible.
func formatStringValue(s string, f Field) (string, error) {
	trimmed := strings.TrimSpace(s)

	switch f.Type {
	case "int":
		if trimmed == "" {
			return "", nil
		}
		val, err := strconv.ParseInt(trimmed, 10, 64)
		if err != nil {
			return "", fmt.Errorf("cannot convert string '%s' to int: %w", s, err)
		}
		return fmt.Sprintf("%d", val), nil

	case "float":
		if trimmed == "" {
			return "", nil
		}
		parsed, err := strconv.ParseFloat(trimmed, 64)
		if err != nil {
			return "", fmt.Errorf("cannot convert string '%s' to float: %w", s, err)
		}
		if f.Decimal > 0 {
			mult := math.Pow(10, float64(f.Decimal))
			val := int64(math.Round(parsed * mult))
			return fmt.Sprintf("%d", val), nil
		}
		return strconv.FormatFloat(parsed, 'f', -1, 64), nil

	case "date":
		if trimmed == "" {
			return "", nil
		}
		format := f.Format
		if format == "" {
			format = "20060102"
		}
		parsed, err := time.Parse(format, trimmed)
		if err != nil {
			return "", fmt.Errorf("cannot convert string '%s' to date with format '%s': %w", s, format, err)
		}
		return parsed.Format(format), nil
	}

	return s, nil
}
