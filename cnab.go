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

type Encoder struct{}

func NewEncoder() *Encoder {
	return &Encoder{}
}

func (e *Encoder) Encode(v interface{}) ([]byte, error) {
	// Implementation to be added in encode.go
	return encode(v)
}

type Decoder struct{}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (d *Decoder) Decode(data []byte, v interface{}) error {
	// Implementation to be added in decode.go
	return decode(data, v)
}
