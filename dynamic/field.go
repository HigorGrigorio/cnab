package dynamic

// Field describes a CNAB field definition used for dynamic layouts.
type Field struct {
	Name     string `json:"name"`
	Size     int    `json:"size"`
	Start    int    `json:"start,omitempty"` // 1-based start position
	Required bool   `json:"required,omitempty"`

	// Formatting options
	Fill  string `json:"fill,omitempty"`  // Character to fill with (default ' ' or '0')
	Align string `json:"align,omitempty"` // "left" or "right"

	// Type specific
	Type    string `json:"type,omitempty"`    // "string", "int", "float", "date"
	Format  string `json:"format,omitempty"`  // Date format
	Decimal int    `json:"decimal,omitempty"` // Decimal places for float/int
}

// Fields represents a collection of CNAB field definitions.
type Fields []Field
