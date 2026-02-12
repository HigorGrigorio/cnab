package cnab

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
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
	
	if ft.end != 0 {
		if ft.start == 0 {
			return ft, errors.Wrap(ErrInvalidTag, "start is mandatory when end is set")
		}

		if ft.start > ft.end {
			return ft, errors.Wrap(ErrInvalidTag, fmt.Sprintf("start %d must be less than end %d", ft.start, ft.end))
		}
		
		// if size is set, it must match the interval
		if ft.size != 0 && ft.size != ft.end-ft.start+1 {
			return ft, errors.Wrap(ErrInvalidTag, fmt.Sprintf("size %d does not match interval %d-%d", ft.size, ft.start, ft.end))
		}
		
		ft.size = ft.end - ft.start + 1
	} else {
		// auto size
		if ft.size == 0 {
			if ft.start == 0 {
				// for literal, size is the length of the literal
				if ft.literalValue != "" {
					ft.size = len(ft.literalValue)
				} else if ft.format != "" {
					ft.size = len(ft.format)
				}
			} else {
				// in this case has start and size
				ft.end = ft.start + len(ft.literalValue) - 1
				ft.size = ft.end - ft.start + 1
			}

			if ft.size == 0 {
				return ft, errors.Wrap(ErrInvalidTag, "size is mandatory when end is set")
			}
		}
	}

	return ft, nil
}
