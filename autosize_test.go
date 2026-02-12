package cnab

import (
	"testing"
)

func TestAutosizeWithoutEnd(t *testing.T) {
	type AutoStruct struct {
		Prefix string `cnab:"literal:PRE"`          // Expected size: 3
		Data   string `cnab:"size:5;literal:DATA_"` // Explicit size matching literal
		Suffix string `cnab:"literal:SUF"`          // Expected size: 3
	}

	s := AutoStruct{}
	// PRE (3) + DATA_ (5) + SUF (3) = PREDATA_SUF
	expected := "PREDATA_SUF"

	data, err := Marshal(s)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	if string(data) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, string(data))
	}
}

func TestAutoSizeWithSize(t *testing.T) {
	type PredefiniedSize struct {
		Prefix string `cnab:"literal:PRE"`
		Data   string `cnab:"size:5;literal:H;align:left"`
		Suffix string `cnab:"literal:SUF"`
	}

	s := PredefiniedSize{}
	expected := "PREH    SUF"

	data, err := Marshal(s)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	if string(data) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, string(data))
	}
}

func TestAutosizeWithEnd(t *testing.T) {
	type AutoEndStruct struct {
		// End: 5. Literal: HELLO (5 chars). Start should be 1.
		Field1 string `cnab:"start:1;literal:HELLO"`
	}

	s := AutoEndStruct{}
	expected := "HELLO"

	data, err := Marshal(s)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	if string(data) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, string(data))
	}
}
