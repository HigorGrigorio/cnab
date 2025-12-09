package cnab

import "errors"

var (
	ErrInvalidStruct       = errors.New("cnab: interface must be a struct")
	ErrInvalidStructPtr    = errors.New("cnab: interface must be a pointer to struct")
	ErrInvalidTag          = errors.New("cnab: invalid tag format")
	ErrLineTooShort        = errors.New("cnab: line is too short for defined fields")
	ErrFieldSizeMismatch   = errors.New("cnab: field value size mismatch")
	ErrUnsupportedType     = errors.New("cnab: unsupported field type")
	ErrInvalidNumberFormat = errors.New("cnab: invalid number format")
	ErrInvalidDateFormat   = errors.New("cnab: invalid date format")
)

