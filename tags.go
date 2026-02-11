package cnab

import (
	"strconv"
	"strings"
)

type fieldTag struct {
	start        int
	end          int
	size         int
	fill         rune
	align        string // "left" or "right"
	format       string
	decimal      int
	hasDate      bool
	literalValue string
}

func parseTag(tag string) (fieldTag, error) {
	ft := fieldTag{
		fill:  ' ',
		align: "left", // Default alignment
	}

	parts := strings.Split(tag, ";")
	for _, part := range parts {
		kv := strings.SplitN(part, ":", 2)
		key := strings.TrimSpace(kv[0])
		value := ""
		if len(kv) > 1 {
			value = strings.TrimSpace(kv[1])
		}

		switch key {
		case "start":
			v, err := strconv.Atoi(value)
			if err != nil {
				return ft, ErrInvalidTag
			}
			ft.start = v
		case "size":
			v, err := strconv.Atoi(value)
			if err != nil {
				return ft, ErrInvalidTag
			}
			ft.size = v
		case "end":
			v, err := strconv.Atoi(value)
			if err != nil {
				return ft, ErrInvalidTag
			}
			ft.end = v
		case "fill":
			if len(value) == 0 {
				ft.fill = ' '
			} else if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") && len(value) == 3 {
				// Handle '0', ' ' etc.
				ft.fill = rune(value[1])
			} else {
				// Fallback or simple char
				runes := []rune(value)
				if len(runes) > 0 {
					ft.fill = runes[0]
				}
			}
		case "align":
			if value == "right" {
				ft.align = "right"
			} else {
				ft.align = "left"
			}
		case "format":
			ft.format = value
			ft.hasDate = true
		case "decimal":
			v, err := strconv.Atoi(value)
			if err != nil {
				return ft, ErrInvalidTag
			}
			ft.decimal = v
		case "cnab":
			// ignore, just the tag name itself if passed incorrectly
		case "literal":
			ft.literalValue = value
		default:
			// ignore unknown keys or handle as error?
			// For now, ignore to be flexible
		}
	}

	// 1. If size is missing but we have a literal, autosize it.
	if ft.size == 0 && ft.literalValue != "" {
		ft.size = len(ft.literalValue)
	}

	// 2. Derive size from start/end if provided
	if ft.end > 0 && ft.start > 0 {
		if ft.end < ft.start {
			return ft, ErrInvalidTag // invalid interval
		}
		derivedSize := ft.end - ft.start + 1
		if ft.size > 0 && ft.size != derivedSize {
			// specific size conflicts with interval
			return ft, ErrInvalidTag
		}
		ft.size = derivedSize
	}

	// 3. If we have end and size but no start, derive start
	if ft.end > 0 && ft.size > 0 && ft.start == 0 {
		ft.start = ft.end - ft.size + 1
	}

	// 4. If we have start and size but no end, derive end (optional, but good for completeness)
	if ft.start > 0 && ft.size > 0 && ft.end == 0 {
		ft.end = ft.start + ft.size - 1
	}

	if ft.size <= 0 {
		return ft, ErrInvalidTag // Size is mandatory (either explicit or derived)
	}

	return ft, nil
}
