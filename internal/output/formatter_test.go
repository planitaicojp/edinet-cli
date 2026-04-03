package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

type testRow struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestTableFormatter(t *testing.T) {
	var buf bytes.Buffer
	f := New("table")
	rows := []testRow{
		{ID: "1", Name: "Alice", Age: 30},
		{ID: "2", Name: "Bob", Age: 25},
	}
	if err := f.Format(&buf, rows); err != nil {
		t.Fatalf("format error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "ID") || !strings.Contains(out, "NAME") {
		t.Errorf("expected headers, got:\n%s", out)
	}
	if !strings.Contains(out, "Alice") || !strings.Contains(out, "Bob") {
		t.Errorf("expected data, got:\n%s", out)
	}
}

func TestJSONFormatter(t *testing.T) {
	var buf bytes.Buffer
	f := New("json")
	row := testRow{ID: "1", Name: "Alice", Age: 30}
	if err := f.Format(&buf, row); err != nil {
		t.Fatalf("format error: %v", err)
	}
	var result testRow
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}
	if result.Name != "Alice" {
		t.Errorf("expected Alice, got %s", result.Name)
	}
}

func TestCSVFormatter(t *testing.T) {
	var buf bytes.Buffer
	f := New("csv")
	rows := []testRow{
		{ID: "1", Name: "Alice", Age: 30},
		{ID: "2", Name: "Bob", Age: 25},
	}
	if err := f.Format(&buf, rows); err != nil {
		t.Fatalf("format error: %v", err)
	}
	out := buf.String()
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines (header + 2 rows), got %d", len(lines))
	}
	if !strings.Contains(lines[0], "id") {
		t.Errorf("expected csv header with 'id', got: %s", lines[0])
	}
}

func TestNewDefault(t *testing.T) {
	f := New("")
	if _, ok := f.(*TableFormatter); !ok {
		t.Errorf("expected TableFormatter as default, got %T", f)
	}
}
