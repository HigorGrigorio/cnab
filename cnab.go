package cnab

// Marshaler is the interface implemented by types that
// can marshal themselves into a CNAB format.
type Marshaler interface {
	MarshalCNAB() ([]byte, error)
}

// Unmarshaler is the interface implemented by types that
// can unmarshal a CNAB description of themselves.
type Unmarshaler interface {
	UnmarshalCNAB([]byte) error
}

// Marshal returns the CNAB encoding of v.
func Marshal(v interface{}) ([]byte, error) {
	return NewEncoder().Encode(v)
}

// Unmarshal parses the CNAB-encoded data and stores the result
// in the value pointed to by v.
func Unmarshal(data []byte, v interface{}) error {
	return NewDecoder().Decode(data, v)
}

// Encoder provides CNAB encoding for struct values using field tags.
type Encoder struct{}

// NewEncoder creates a new CNAB encoder with default settings.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encode marshals the provided value into CNAB-formatted bytes.
func (e *Encoder) Encode(v interface{}) ([]byte, error) {
	return encode(v)
}

// Decoder provides CNAB decoding for tagged struct values.
type Decoder struct{}

// NewDecoder creates a new CNAB decoder with default settings.
func NewDecoder() *Decoder {
	return &Decoder{}
}

// Decode parses CNAB-formatted data into the provided destination value.
func (d *Decoder) Decode(data []byte, v interface{}) error {
	return decode(data, v)
}
