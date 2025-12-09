package dynamic

import (
	"strings"
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
	layout := []Field{
		{Name: "ID", Size: 5, Type: "int"}, // Defaults: Fill=0, Align=right
		{Name: "Desc", Size: 10},           // Defaults: Fill=" ", Align=left
	}

	data := map[string]interface{}{
		"ID":   "50", // coercion from string to int
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

func TestMarshalStringCoercionSuccess(t *testing.T) {
	layout := []Field{
		{Name: "Code", Size: 3, Type: "int", Fill: "0", Align: "right"},
		{Name: "Name", Size: 10, Type: "string", Fill: " ", Align: "left"},
		{Name: "Date", Size: 8, Type: "date", Format: "20060102"},
		{Name: "Amount", Size: 10, Type: "float", Decimal: 2, Fill: "0", Align: "right"},
	}

	data := map[string]interface{}{
		"Code":   "1",
		"Name":   "TEST",
		"Date":   "20231025",
		"Amount": "123.45",
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

func TestMarshalStringFloatRounding(t *testing.T) {
	layout := []Field{
		{Name: "Amount", Size: 5, Type: "float", Decimal: 2, Fill: "0", Align: "right"},
	}

	data := map[string]interface{}{
		"Amount": "12.345", // rounds to 12.35 -> 1235
	}

	expected := "01235"
	res, err := Marshal(data, layout)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	if string(res) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, string(res))
	}
}

func TestMarshalStringIntError(t *testing.T) {
	layout := []Field{
		{Name: "Code", Size: 3, Type: "int", Fill: "0", Align: "right"},
	}

	data := map[string]interface{}{
		"Code": "ABC",
	}

	_, err := Marshal(data, layout)
	if err == nil {
		t.Fatalf("expected error for invalid int coercion, got nil")
	}

	if !strings.Contains(err.Error(), "cannot convert string 'ABC' to int") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestMarshalStringDateError(t *testing.T) {
	layout := []Field{
		{Name: "Date", Size: 8, Type: "date", Format: "20060102"},
	}

	data := map[string]interface{}{
		"Date": "25/10/2023",
	}

	_, err := Marshal(data, layout)
	if err == nil {
		t.Fatalf("expected error for invalid date coercion, got nil")
	}

	if !strings.Contains(err.Error(), "cannot convert string '25/10/2023' to date") {
		t.Fatalf("unexpected error message: %v", err)
	}
}
