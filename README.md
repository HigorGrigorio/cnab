# CNAB Library for Go

This library provides a flexible engine for parsing and generating fixed-width (positional) text files, commonly used in banking and legacy systems (CNAB).

## Features

- **Struct Tags**: Define layout using `cnab:"..."` tags.
- **Marshal/Unmarshal**: Standard Go interface for encoding/decoding.
- **Dynamic Layouts**: Generate CNAB from `map` and runtime field definitions (ideal for CSV -> CNAB).
- **Support for**:
  - **Strings**: Aligned left/right with custom padding.
  - **Integers**: Supports `int`, `int8-64`, `uint`, `uint8-64`. Handles negative numbers with zero padding correctly (e.g. `-1` -> `-0001` with size 5).
  - **Floats**: Implicit decimal point handling (conversion from integer representation in file to float in struct).
  - **Dates**: Custom formats (e.g. `YYYYMMDD`).
  - **Custom Types**: Implement `Marshaler` and `Unmarshaler` interfaces for full control over field encoding/decoding.

## Example

```go
// Example:                                                                 |20251208|       JOHN|0000010302|000001|
struct Line { //                                                                 |        |           |         |
    Date Date `cnab:"date;format:AAAAMMDD;required"` //  <-----------------------|        |           |         |
    Name string `cnab:"name;required;size:10;align:right;ascii:' '"` // <-----------------|           |         |
    Amount float64 `cnab:"decimal:2;required;size:10"` // <-------------------------------------------|         |
    Counter int `cnab:"counter;required;size:6;align:right;ascii:'0'"` // <-------------------------------------|
}
```

## Usage

### Defining a Layout

```go
type Header struct {
    Code   int       `cnab:"size:3;fill:0;align:right"`
    Name   string    `cnab:"size:10;fill: ;align:left"`
    Date   time.Time `cnab:"size:8;format:20060102"`
    Amount float64   `cnab:"size:10;decimal:2;fill:0;align:right"`
    Type   string    `cnab:"literal:REM;size:3"` // Explicit size
    Ver    string    `cnab:"literal:040"`        // Autosize to 3
}
```

### Marshal (Writing)

```go
h := Header{
    Code:   1,
    Name:   "TEST",
    Date:   time.Now(),
    Amount: 123.45,
}
data, err := cnab.Marshal(h)
// Result: "001TEST      202312080000012345"
```

### Unmarshal (Reading)

```go
line := "001TEST      202312080000012345"
var h Header
err := cnab.Unmarshal([]byte(line), &h)
```

### Dynamic Layout (CSV -> CNAB)

You can generate CNAB lines from a map (e.g., parsed from CSV) without defining a struct:

```go
import "github.com/HigorGrigorio/cnab/dynamic"

layout := []dynamic.Field{
    {Name: "ID", Size: 5, Type: "int"}, // Defaults: Fill="0", Align="right"
    {Name: "Desc", Size: 10, Fill: " ", Align: "left"},
}

data := map[string]interface{}{
    "ID": 50,
    "Desc": "ITEM",
}

line, err := dynamic.Marshal(data, layout)
// Result: "00050ITEM      "
```

String inputs are coerced automatically when `Type` is set:
- `int`/`float`: strings são convertidas; para `Decimal`, valores com ponto são aceitos e arredondados antes de aplicar o padding. Erros trazem a string original e a causa do parse.
- `date`: converte usando `Format` (ou `20060102` se vazio) e falha com mensagem clara quando o texto não obedece ao formato.

### Custom Encoding/Decoding

You can implement `MarshalCNAB` and `UnmarshalCNAB` for custom types:

```go
type BoolFlag bool

func (b BoolFlag) MarshalCNAB() ([]byte, error) {
    if b {
        return []byte("S"), nil
    }
    return []byte("N"), nil
}

func (b *BoolFlag) UnmarshalCNAB(data []byte) error {
    s := string(data)
    if s == "S" {
        *b = true
    } else {
        *b = false
    }
    return nil
}

type MyStruct struct {
    Flag BoolFlag `cnab:"size:1"`
}
```

## Supported Tags (struct)

| Tag     | Description                                                | Default                                 | Notes / Rules                                                                                |
|---------|------------------------------------------------------------|-----------------------------------------|----------------------------------------------------------------------------------------------|
| `size`  | Field length in characters. Required unless `end` is set. | –                                       | If `end` is provided, `size` is derived as `end - start + 1`.                                |
| `start` | 1-based start position. If omitted, continues sequentially.| –                                       | Used with `size` or `end` to define the interval.                                            |
| `end`   | 1-based end position.                                      | –                                       | When present, it takes precedence over `size` (interval = `start..end`).                     |
| `fill`  | Padding character.                                         | `' '` for strings; `'0'` for numbers    | Applied to reach `size`.                                                                    |
| `align` | Padding direction (`left` or `right`).                     | `left` for strings; `right` for numbers | Works with `fill` to place the value within the field.                                       |
| `format`| Date format for `time.Time` fields.                        | `20060102`                              | Go time layout.                                                                              |
| `decimal`| Implied decimal places for numeric fields.                | 0                                       | Value is multiplied/divided by `10^decimal` on marshal/unmarshal.                            |
| `literal`| Constant value override.                                   | –                                       | Always outputs this value. Used for autosize if `size` is missing.                           |

Positioning rules: 
- If `start` is omitted, the field begins right after the previous one. 
- If `size` is omitted but `literal` is present, `size` defaults to `len(literal)`.
- If `start` is omitted but `end` and `size` (explicit or derived) are present, `start` is calculated as `end - size + 1`.
- Overlapping intervals cause an encode error. 
- Negative numbers keep the sign on the left with zero-fill (e.g., `-1` in size 5 becomes `-0001`).
