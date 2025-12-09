package cnab

import "errors"

var (
	// ErrInvalidStruct indicates that a non-struct value was provided.
	ErrInvalidStruct = errors.New("cnab: interface must be a struct")

	// ErrInvalidStructPtr indicates that a non-pointer or nil pointer was provided.
	ErrInvalidStructPtr = errors.New("cnab: interface must be a pointer to struct")

	// ErrInvalidTag indicates that a struct field tag could not be parsed.
	ErrInvalidTag = errors.New("cnab: invalid tag format")

	// ErrLineTooShort indicates that the input line does not contain all tagged fields.
	ErrLineTooShort = errors.New("cnab: line is too short for defined fields")

	// ErrFieldSizeMismatch indicates that a value does not fit in its declared size.
	ErrFieldSizeMismatch = errors.New("cnab: field value size mismatch")

	// ErrUnsupportedType indicates that a field has a type that cannot be encoded or decoded.
	ErrUnsupportedType = errors.New("cnab: unsupported field type")

	// ErrInvalidNumberFormat indicates that a numeric field could not be parsed.
	ErrInvalidNumberFormat = errors.New("cnab: invalid number format")

	// ErrInvalidDateFormat indicates that a date field could not be parsed with the provided format.
	ErrInvalidDateFormat = errors.New("cnab: invalid date format")
)
