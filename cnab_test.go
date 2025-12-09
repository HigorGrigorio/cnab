package cnab

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

type TestHeader struct {
	Code      int       `cnab:"size:3;fill:0;align:right"`
	Name      string    `cnab:"size:10;fill: ;align:left"`
	Date      time.Time `cnab:"size:8;format:20060102"`
	Amount    float64   `cnab:"size:10;decimal:2;fill:0;align:right"`
	EmptyDate time.Time `cnab:"size:8;format:20060102"`
}

func TestMarshal(t *testing.T) {
	dt := time.Date(2023, 10, 25, 0, 0, 0, 0, time.UTC)
	h := TestHeader{
		Code:   1,
		Name:   "TEST",
		Date:   dt,
		Amount: 123.45,
	}

	expected := "001TEST      202310250000012345        "
	data, err := Marshal(h)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	if string(data) != expected {
		t.Errorf("Marshal mismatch:\nGot:  '%s'\nWant: '%s'", string(data), expected)
	}
}

func TestUnmarshal(t *testing.T) {
	input := "001TEST      202310250000012345        "
	var h TestHeader

	err := Unmarshal([]byte(input), &h)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if h.Code != 1 {
		t.Errorf("Expected Code 1, got %d", h.Code)
	}
	if h.Name != "TEST" {
		t.Errorf("Expected Name 'TEST', got '%s'", h.Name)
	}

	expectedDate := time.Date(2023, 10, 25, 0, 0, 0, 0, time.UTC)
	if !h.Date.Equal(expectedDate) {
		t.Errorf("Expected Date %v, got %v", expectedDate, h.Date)
	}

	if h.Amount != 123.45 {
		t.Errorf("Expected Amount 123.45, got %f", h.Amount)
	}

	if !h.EmptyDate.IsZero() {
		t.Errorf("Expected EmptyDate to be zero, got %v", h.EmptyDate)
	}
}

func TestMarshalOverflow(t *testing.T) {
	h := TestHeader{
		Name: "THIS NAME IS TOO LONG",
	}
	_, err := Marshal(h)
	if err == nil {
		t.Error("Expected error for overflow, got nil")
	}
}

func TestIntegerValues(t *testing.T) {
	type IntStruct struct {
		Val     int   `cnab:"size:5;fill:0;align:right"`
		NegVal  int   `cnab:"size:5;fill:0;align:right"`
		BigVal  int64 `cnab:"size:10;fill:0;align:right"`
		UintVal uint  `cnab:"size:5;fill:0;align:right"`
	}

	s := IntStruct{
		Val:     123,
		NegVal:  -1,
		BigVal:  1234567890,
		UintVal: 99,
	}

	// 00123
	// -0001 (5 chars)
	// 1234567890
	// 00099

	expected := "00123-0001123456789000099"
	data, err := Marshal(s)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	if string(data) != expected {
		t.Errorf("Marshal IntStruct mismatch:\nGot:  '%s'\nWant: '%s'", string(data), expected)
	}

	var decoded IntStruct
	err = Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Unmarshal IntStruct failed: %v", err)
	}

	if decoded.Val != 123 {
		t.Errorf("Expected Val 123, got %d", decoded.Val)
	}
	if decoded.NegVal != -1 {
		t.Errorf("Expected NegVal -1, got %d", decoded.NegVal)
	}
	if decoded.BigVal != 1234567890 {
		t.Errorf("Expected BigVal 1234567890, got %d", decoded.BigVal)
	}
	if decoded.UintVal != 99 {
		t.Errorf("Expected UintVal 99, got %d", decoded.UintVal)
	}
}

// CustomType implements Marshaler and Unmarshaler
type CustomBool bool

func (c CustomBool) MarshalCNAB() ([]byte, error) {
	if c {
		return []byte("S"), nil
	}
	return []byte("N"), nil
}

func (c *CustomBool) UnmarshalCNAB(data []byte) error {
	s := strings.TrimSpace(string(data))
	if s == "S" {
		*c = true
	} else if s == "N" {
		*c = false
	} else {
		return fmt.Errorf("invalid boolean value: %s", s)
	}
	return nil
}

func TestCustomInterface(t *testing.T) {
	type CustomStruct struct {
		Flag CustomBool `cnab:"size:1"`
	}

	s := CustomStruct{Flag: true}

	data, err := Marshal(s)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	if string(data) != "S" {
		t.Errorf("Expected 'S', got '%s'", string(data))
	}

	var s2 CustomStruct
	err = Unmarshal(data, &s2)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if s2.Flag != true {
		t.Errorf("Expected true, got false")
	}

	// Test false case
	s3 := CustomStruct{Flag: false}
	data, err = Marshal(s3)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	if string(data) != "N" {
		t.Errorf("Expected 'N', got '%s'", string(data))
	}
}
