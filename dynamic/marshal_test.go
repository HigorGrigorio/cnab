package dynamic

import (
	"testing"
	"time"
)

func TestMarshal(t *testing.T) {
	layout := []Field{
		{Name: "Code", Size: 3, Type: "int", Fill: "0", Align: "right"},
		{Name: "Name", Size: 10, Type: "string", Fill: " ", Align: "left"},
		{Name: "Date", Size: 8, Type: "date", Format: "20060102"},
		{Name: "Amount", Size: 10, Type: "float", Decimal: 2, Fill: "0", Align: "right"},
	}

	dt := time.Date(2023, 10, 25, 0, 0, 0, 0, time.UTC)
	data := map[string]interface{}{
		"Code":   1,
		"Name":   "TEST",
		"Date":   dt,
		"Amount": 123.45,
	}

	expected := "001TEST      202310250000012345"
	res, err := Marshal(data, layout)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	if string(res) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, string(res))
	}
}

func TestMarshalCSVLike(t *testing.T) {
	// Simulating CSV reading where everything might be string initially,
	// but user converts to map[string]interface{} with correct types or we handle string conversion?
	// Our formatValue handles specific types. If user passes "123" as string for an int field,
	// our formatValue will treat it as string "123".
	// If the layout says Type: "int" but value is string "123", we might want to auto-convert?
	// The current impl of formatValue does fmt.Sprintf("%v") for unknown, or specific type switch.
	// If Type="int" and val="123", it falls to default string print, which is fine.
	// But padding defaults depend on Type.

	layout := []Field{
		{Name: "ID", Size: 5, Type: "int"}, // Defaults: Fill=0, Align=right
		{Name: "Desc", Size: 10},           // Defaults: Fill=" ", Align=left
	}

	data := map[string]interface{}{
		"ID":   50,
		"Desc": "ITEM",
	}

	// ID: 00050
	// Desc: ITEM
	expected := "00050ITEM      "
	res, err := Marshal(data, layout)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	if string(res) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, string(res))
	}
}
