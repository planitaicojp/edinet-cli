# edinet-cli Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a Go CLI tool for EDINET API v2 with minimal working functionality, goreleaser for brew/scoop distribution, and README.

**Architecture:** Cobra-based nested subcommands (`edinet document list/download`, `edinet company search`), `internal/` for API client, config, models, output formatters, errors. Follows conoha-cli patterns.

**Tech Stack:** Go, Cobra, goreleaser, Homebrew tap, Scoop bucket

---

## File Map

### Phase 1: Project Scaffold + Config
| File | Responsibility |
|------|----------------|
| `go.mod` | Module definition |
| `main.go` | Entry point, calls `cmd.Execute()` |
| `cmd/root.go` | Root command, global flags (`--format`, `--api-key`, `--verbose`, `--no-color`) |
| `cmd/version.go` | `edinet version` — prints version (injected via ldflags) |
| `internal/config/env.go` | Environment variable constants |
| `internal/config/config.go` | Load/Save `~/.config/edinet/config.yaml` |
| `internal/errors/exitcodes.go` | Exit code constants |
| `internal/errors/errors.go` | Custom error types |

### Phase 2: API Client + Output
| File | Responsibility |
|------|----------------|
| `internal/api/client.go` | HTTP client with retry, debug logging, Subscription-Key injection |
| `internal/model/document.go` | Document, Metadata, ResultSet structs |
| `internal/model/response.go` | DocumentListResponse wrapper |
| `internal/output/formatter.go` | Formatter interface + factory |
| `internal/output/table.go` | Table formatter (tabwriter) |
| `internal/output/json.go` | JSON formatter |
| `internal/output/csv.go` | CSV formatter |

### Phase 3: Document Commands
| File | Responsibility |
|------|----------------|
| `cmd/cmdutil/client.go` | `NewClient(cmd)` factory |
| `cmd/cmdutil/format.go` | `GetFormat(cmd)` helper |
| `cmd/cmdutil/flags.go` | Global flag accessors |
| `cmd/document/document.go` | `edinet document` group command |
| `cmd/document/list.go` | `edinet document list --date YYYY-MM-DD` |
| `internal/api/documents.go` | `ListDocuments()`, `DownloadDocument()` |
| `cmd/document/download.go` | `edinet document download <docID>` |

### Phase 4: Config Command
| File | Responsibility |
|------|----------------|
| `cmd/config/config.go` | `edinet config` group command |
| `cmd/config/set.go` | `edinet config set api-key <value>` |
| `cmd/config/show.go` | `edinet config show` |

### Phase 5: Release + README + Makefile
| File | Responsibility |
|------|----------------|
| `.goreleaser.yaml` | Build/release config with brew+scoop |
| `.golangci.yml` | golangci-lint configuration |
| `Makefile` | build, test, lint, clean targets |
| `.gitignore` | Go + build artifacts |
| `LICENSE` | Apache-2.0 |
| `README.md` | Installation, usage, examples |

---

## Task 1: Go Module + Entry Point

**Files:**
- Create: `go.mod`
- Create: `main.go`

- [ ] **Step 1: Initialize Go module**

```bash
cd /root/dev/planitai/edinet-cli
go mod init github.com/planitaicojp/edinet-cli
```

- [ ] **Step 2: Create main.go**

```go
package main

import (
	"fmt"
	"os"

	"github.com/planitaicojp/edinet-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```

- [ ] **Step 3: Verify it compiles (will fail — cmd package doesn't exist yet)**

```bash
go build ./...
```

Expected: error about `cmd` package not found. This confirms main.go is syntactically correct.

---

## Task 2: Exit Codes + Error Types

**Files:**
- Create: `internal/errors/exitcodes.go`
- Create: `internal/errors/errors.go`
- Create: `internal/errors/errors_test.go`

- [ ] **Step 1: Write error tests**

```go
// internal/errors/errors_test.go
package errors

import (
	"testing"
)

func TestAPIError(t *testing.T) {
	err := &APIError{StatusCode: 400, Code: "BAD_REQUEST", Message: "invalid date"}
	if err.Error() != "API error (400): BAD_REQUEST - invalid date" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
	if err.ExitCode() != ExitAPI {
		t.Errorf("expected exit code %d, got %d", ExitAPI, err.ExitCode())
	}
}

func TestAuthError(t *testing.T) {
	err := &AuthError{Message: "invalid API key"}
	if err.ExitCode() != ExitAuth {
		t.Errorf("expected exit code %d, got %d", ExitAuth, err.ExitCode())
	}
}

func TestNotFoundError(t *testing.T) {
	err := &NotFoundError{Resource: "document", ID: "S1234567"}
	if err.Error() != "document not found: S1234567" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
	if err.ExitCode() != ExitNotFound {
		t.Errorf("expected exit code %d, got %d", ExitNotFound, err.ExitCode())
	}
}

func TestValidationError(t *testing.T) {
	err := &ValidationError{Field: "date", Message: "invalid format"}
	if err.ExitCode() != ExitValidation {
		t.Errorf("expected exit code %d, got %d", ExitValidation, err.ExitCode())
	}
}

func TestNetworkError(t *testing.T) {
	err := &NetworkError{Err: fmt.Errorf("connection refused")}
	if err.ExitCode() != ExitNetwork {
		t.Errorf("expected exit code %d, got %d", ExitNetwork, err.ExitCode())
	}
}

func TestGetExitCode(t *testing.T) {
	tests := []struct {
		err      error
		expected int
	}{
		{&AuthError{Message: "no key"}, ExitAuth},
		{&APIError{StatusCode: 500}, ExitAPI},
		{fmt.Errorf("unknown"), ExitGeneral},
	}
	for _, tt := range tests {
		if got := GetExitCode(tt.err); got != tt.expected {
			t.Errorf("GetExitCode(%v) = %d, want %d", tt.err, got, tt.expected)
		}
	}
}
```

Add `"fmt"` import at the top of the test file.

- [ ] **Step 2: Run tests to verify they fail**

```bash
cd /root/dev/planitai/edinet-cli
go test ./internal/errors/...
```

Expected: FAIL — types not defined yet.

- [ ] **Step 3: Create exitcodes.go**

```go
// internal/errors/exitcodes.go
package errors

const (
	ExitOK         = 0
	ExitGeneral    = 1
	ExitValidation = 2
	ExitAuth       = 3
	ExitAPI        = 4
	ExitNetwork    = 5
	ExitNotFound   = 6
)

// ExitCoder is implemented by errors that have a specific exit code.
type ExitCoder interface {
	ExitCode() int
}

// GetExitCode returns the exit code for an error.
func GetExitCode(err error) int {
	if ec, ok := err.(ExitCoder); ok {
		return ec.ExitCode()
	}
	return ExitGeneral
}
```

- [ ] **Step 4: Create errors.go**

```go
// internal/errors/errors.go
package errors

import "fmt"

type APIError struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (%d): %s - %s", e.StatusCode, e.Code, e.Message)
}

func (e *APIError) ExitCode() int { return ExitAPI }

type AuthError struct {
	Message string
}

func (e *AuthError) Error() string   { return fmt.Sprintf("authentication error: %s", e.Message) }
func (e *AuthError) ExitCode() int   { return ExitAuth }

type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string   { return fmt.Sprintf("%s not found: %s", e.Resource, e.ID) }
func (e *NotFoundError) ExitCode() int   { return ExitNotFound }

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string   { return fmt.Sprintf("validation error: %s: %s", e.Field, e.Message) }
func (e *ValidationError) ExitCode() int   { return ExitValidation }

type NetworkError struct {
	Err error
}

func (e *NetworkError) Error() string   { return fmt.Sprintf("network error: %s", e.Err) }
func (e *NetworkError) Unwrap() error   { return e.Err }
func (e *NetworkError) ExitCode() int   { return ExitNetwork }
```

- [ ] **Step 5: Run tests to verify they pass**

```bash
go test ./internal/errors/... -v
```

Expected: All 6 tests PASS.

- [ ] **Step 6: Commit**

```bash
git add internal/errors/
git commit -m "feat: add custom error types and exit codes"
```

---

## Task 3: Config (env + yaml)

**Files:**
- Create: `internal/config/env.go`
- Create: `internal/config/config.go`
- Create: `internal/config/config_test.go`

- [ ] **Step 1: Write config tests**

```go
// internal/config/config_test.go
package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEmpty(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(EnvConfigDir, dir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.APIKey != "" {
		t.Errorf("expected empty api key, got %q", cfg.APIKey)
	}
	if cfg.DefaultFormat != "" {
		t.Errorf("expected empty default format, got %q", cfg.DefaultFormat)
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(EnvConfigDir, dir)

	cfg := &Config{
		APIKey:        "test-key-12345",
		DefaultFormat: "json",
	}
	if err := Save(cfg); err != nil {
		t.Fatalf("save error: %v", err)
	}

	// Verify file permissions
	info, err := os.Stat(filepath.Join(dir, "config.yaml"))
	if err != nil {
		t.Fatalf("stat error: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected permissions 0600, got %o", info.Mode().Perm())
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if loaded.APIKey != "test-key-12345" {
		t.Errorf("expected api key 'test-key-12345', got %q", loaded.APIKey)
	}
	if loaded.DefaultFormat != "json" {
		t.Errorf("expected default format 'json', got %q", loaded.DefaultFormat)
	}
}

func TestResolveAPIKey_EnvOverridesConfig(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(EnvConfigDir, dir)
	t.Setenv(EnvAPIKey, "env-key")

	cfg := &Config{APIKey: "config-key"}
	if err := Save(cfg); err != nil {
		t.Fatalf("save error: %v", err)
	}

	got := ResolveAPIKey("")
	if got != "env-key" {
		t.Errorf("expected 'env-key', got %q", got)
	}
}

func TestResolveAPIKey_FlagOverridesEnv(t *testing.T) {
	t.Setenv(EnvAPIKey, "env-key")

	got := ResolveAPIKey("flag-key")
	if got != "flag-key" {
		t.Errorf("expected 'flag-key', got %q", got)
	}
}

func TestResolveAPIKey_FallbackToConfig(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(EnvConfigDir, dir)
	t.Setenv(EnvAPIKey, "")

	cfg := &Config{APIKey: "config-key"}
	if err := Save(cfg); err != nil {
		t.Fatalf("save error: %v", err)
	}

	got := ResolveAPIKey("")
	if got != "config-key" {
		t.Errorf("expected 'config-key', got %q", got)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/config/... -v
```

Expected: FAIL — types not defined.

- [ ] **Step 3: Create env.go**

```go
// internal/config/env.go
package config

const (
	EnvAPIKey     = "EDINET_API_KEY"
	EnvFormat     = "EDINET_FORMAT"
	EnvConfigDir  = "EDINET_CONFIG_DIR"
	EnvDebug      = "EDINET_DEBUG"
)
```

- [ ] **Step 4: Create config.go**

```go
// internal/config/config.go
package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	APIKey        string `yaml:"api_key"`
	DefaultFormat string `yaml:"default_format,omitempty"`
}

func configDir() string {
	if dir := os.Getenv(EnvConfigDir); dir != "" {
		return dir
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "edinet")
}

func configPath() string {
	return filepath.Join(configDir(), "config.yaml")
}

// Load reads config from disk. Returns empty config if file does not exist.
func Load() (*Config, error) {
	cfg := &Config{}
	data, err := os.ReadFile(configPath())
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Save writes config to disk with 0600 permissions.
func Save(cfg *Config) error {
	dir := configDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(), data, 0600)
}

// ResolveAPIKey returns the API key from flag > env > config (in priority order).
func ResolveAPIKey(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	if env := os.Getenv(EnvAPIKey); env != "" {
		return env
	}
	cfg, err := Load()
	if err != nil {
		return ""
	}
	return cfg.APIKey
}

// ConfigDir returns the config directory path (exported for cache location).
func ConfigDir() string {
	return configDir()
}
```

- [ ] **Step 5: Add yaml.v3 dependency**

```bash
go get gopkg.in/yaml.v3
```

- [ ] **Step 6: Run tests to verify they pass**

```bash
go test ./internal/config/... -v
```

Expected: All 5 tests PASS.

- [ ] **Step 7: Commit**

```bash
git add internal/config/ go.mod go.sum
git commit -m "feat: add config management with env var priority"
```

---

## Task 4: Output Formatters

**Files:**
- Create: `internal/output/formatter.go`
- Create: `internal/output/table.go`
- Create: `internal/output/json.go`
- Create: `internal/output/csv.go`
- Create: `internal/output/formatter_test.go`

- [ ] **Step 1: Write formatter tests**

```go
// internal/output/formatter_test.go
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
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/output/... -v
```

Expected: FAIL.

- [ ] **Step 3: Create formatter.go**

```go
// internal/output/formatter.go
package output

import "io"

// Formatter formats data and writes to w.
type Formatter interface {
	Format(w io.Writer, data any) error
}

// New returns a Formatter for the given format string.
func New(format string) Formatter {
	switch format {
	case "json":
		return &JSONFormatter{}
	case "csv":
		return &CSVFormatter{}
	default:
		return &TableFormatter{}
	}
}
```

- [ ] **Step 4: Create json.go**

```go
// internal/output/json.go
package output

import (
	"encoding/json"
	"io"
)

type JSONFormatter struct{}

func (f *JSONFormatter) Format(w io.Writer, data any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}
```

- [ ] **Step 5: Create table.go**

```go
// internal/output/table.go
package output

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"
)

type TableFormatter struct{}

func (f *TableFormatter) Format(w io.Writer, data any) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	defer tw.Flush()

	v := reflect.ValueOf(data)

	// Handle slice
	if v.Kind() == reflect.Slice {
		if v.Len() == 0 {
			fmt.Fprintln(w, "No results.")
			return nil
		}
		elemType := v.Type().Elem()
		headers := structHeaders(elemType)
		fmt.Fprintln(tw, strings.Join(headers, "\t"))

		for i := 0; i < v.Len(); i++ {
			vals := structValues(v.Index(i))
			fmt.Fprintln(tw, strings.Join(vals, "\t"))
		}
		return nil
	}

	// Handle single struct
	if v.Kind() == reflect.Struct {
		headers := structHeaders(v.Type())
		fmt.Fprintln(tw, strings.Join(headers, "\t"))
		vals := structValues(v)
		fmt.Fprintln(tw, strings.Join(vals, "\t"))
		return nil
	}

	// Fallback
	_, err := fmt.Fprintln(w, data)
	return err
}

func structHeaders(t reflect.Type) []string {
	var headers []string
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		name := strings.Split(tag, ",")[0]
		headers = append(headers, strings.ToUpper(name))
	}
	return headers
}

func structValues(v reflect.Value) []string {
	var vals []string
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		vals = append(vals, fmt.Sprintf("%v", v.Field(i).Interface()))
	}
	return vals
}
```

- [ ] **Step 6: Create csv.go**

```go
// internal/output/csv.go
package output

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strings"
)

type CSVFormatter struct{}

func (f *CSVFormatter) Format(w io.Writer, data any) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Slice {
		if v.Len() == 0 {
			return nil
		}
		elemType := v.Type().Elem()
		headers := csvHeaders(elemType)
		if err := cw.Write(headers); err != nil {
			return err
		}
		for i := 0; i < v.Len(); i++ {
			vals := csvValues(v.Index(i))
			if err := cw.Write(vals); err != nil {
				return err
			}
		}
		return nil
	}

	if v.Kind() == reflect.Struct {
		headers := csvHeaders(v.Type())
		if err := cw.Write(headers); err != nil {
			return err
		}
		return cw.Write(csvValues(v))
	}

	return cw.Write([]string{fmt.Sprint(data)})
}

func csvHeaders(t reflect.Type) []string {
	var headers []string
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		headers = append(headers, strings.Split(tag, ",")[0])
	}
	return headers
}

func csvValues(v reflect.Value) []string {
	var vals []string
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		vals = append(vals, fmt.Sprintf("%v", v.Field(i).Interface()))
	}
	return vals
}
```

- [ ] **Step 7: Run tests to verify they pass**

```bash
go test ./internal/output/... -v
```

Expected: All 4 tests PASS.

- [ ] **Step 8: Commit**

```bash
git add internal/output/
git commit -m "feat: add output formatters (table, json, csv)"
```

---

## Task 5: API Client

**Files:**
- Create: `internal/api/client.go`
- Create: `internal/api/client_test.go`
- Create: `internal/model/document.go`
- Create: `internal/model/response.go`

- [ ] **Step 1: Create model types**

```go
// internal/model/document.go
package model

// Document represents a single filing from the EDINET document list API.
type Document struct {
	SeqNumber            int     `json:"seqNumber"`
	DocID                string  `json:"docID"`
	EdinetCode           string  `json:"edinetCode"`
	SecCode              string  `json:"secCode"`
	JCN                  string  `json:"JCN"`
	FilerName            string  `json:"filerName"`
	FundCode             string  `json:"fundCode"`
	OrdinanceCode        string  `json:"ordinanceCode"`
	FormCode             string  `json:"formCode"`
	DocTypeCode          string  `json:"docTypeCode"`
	PeriodStart          *string `json:"periodStart"`
	PeriodEnd            *string `json:"periodEnd"`
	SubmitDateTime       string  `json:"submitDateTime"`
	DocDescription       string  `json:"docDescription"`
	IssuerEdinetCode     *string `json:"issuerEdinetCode"`
	SubjectEdinetCode    *string `json:"subjectEdinetCode"`
	SubsidiaryEdinetCode *string `json:"subsidiaryEdinetCode"`
	CurrentReportReason  *string `json:"currentReportReason"`
	ParentDocID          *string `json:"parentDocID"`
	OpeDateTime          *string `json:"opeDateTime"`
	WithdrawalStatus     string  `json:"withdrawalStatus"`
	DocInfoEditStatus    string  `json:"docInfoEditStatus"`
	DisclosureStatus     string  `json:"disclosureStatus"`
	XbrlFlag             string  `json:"xbrlFlag"`
	PdfFlag              string  `json:"pdfFlag"`
	AttachDocFlag        string  `json:"attachDocFlag"`
	EnglishDocFlag       string  `json:"englishDocFlag"`
	CsvFlag              string  `json:"csvFlag"`
	LegalStatus          string  `json:"legalStatus"`
}
```

```go
// internal/model/response.go
package model

// Metadata contains API response metadata.
type Metadata struct {
	Title     string          `json:"title"`
	Parameter MetadataParam   `json:"parameter"`
	ResultSet MetadataResult  `json:"resultset"`
	ProcessDateTime string   `json:"processDateTime"`
	Status    string          `json:"status"`
	Message   string          `json:"message"`
}

type MetadataParam struct {
	Date string `json:"date"`
	Type string `json:"type"`
}

type MetadataResult struct {
	Count int `json:"count"`
}

// DocumentListResponse is the response from the document list API.
type DocumentListResponse struct {
	Metadata Metadata   `json:"metadata"`
	Results  []Document `json:"results"`
}
```

- [ ] **Step 2: Write API client tests**

```go
// internal/api/client_test.go
package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_SubscriptionKey(t *testing.T) {
	var gotKey string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotKey = r.URL.Query().Get("Subscription-Key")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(`{"metadata":{"title":"test","parameter":{"date":"2024-01-01","type":"1"},"resultset":{"count":0},"processDateTime":"2024-01-01 00:00","status":"200","message":"OK"}}`))
	}))
	defer ts.Close()

	c := &Client{
		BaseURL: ts.URL,
		APIKey:  "test-api-key",
		HTTP:    ts.Client(),
	}

	_, err := c.Get("/documents.json", map[string]string{"date": "2024-01-01"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotKey != "test-api-key" {
		t.Errorf("expected Subscription-Key 'test-api-key', got %q", gotKey)
	}
}

func TestClient_ErrorHandling(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"metadata":{"status":"403","message":"Invalid API Key"}}`))
	}))
	defer ts.Close()

	c := &Client{
		BaseURL: ts.URL,
		APIKey:  "bad-key",
		HTTP:    ts.Client(),
	}

	_, err := c.Get("/documents.json", map[string]string{"date": "2024-01-01"})
	if err == nil {
		t.Fatal("expected error for 403 response")
	}
}
```

- [ ] **Step 3: Run tests to verify they fail**

```bash
go test ./internal/api/... -v
```

Expected: FAIL.

- [ ] **Step 4: Create client.go**

```go
// internal/api/client.go
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	cerrors "github.com/planitaicojp/edinet-cli/internal/errors"
)

const DefaultBaseURL = "https://api.edinet-fsa.go.jp/api/v2"

type Client struct {
	BaseURL string
	APIKey  string
	HTTP    *http.Client
	Debug   bool
}

func NewClient(apiKey string) *Client {
	return &Client{
		BaseURL: DefaultBaseURL,
		APIKey:  apiKey,
		HTTP: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Get sends a GET request and returns the response body bytes.
func (c *Client) Get(path string, params map[string]string) ([]byte, error) {
	u, err := url.Parse(c.BaseURL + path)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("Subscription-Key", c.APIKey)
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	resp, err := c.HTTP.Get(u.String())
	if err != nil {
		return nil, &cerrors.NetworkError{Err: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, classifyError(resp.StatusCode, body)
	}

	return body, nil
}

// GetBinary sends a GET request and returns the response body as a reader.
func (c *Client) GetBinary(path string, params map[string]string) (io.ReadCloser, string, error) {
	u, err := url.Parse(c.BaseURL + path)
	if err != nil {
		return nil, "", err
	}

	q := u.Query()
	q.Set("Subscription-Key", c.APIKey)
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	resp, err := c.HTTP.Get(u.String())
	if err != nil {
		return nil, "", &cerrors.NetworkError{Err: err}
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, "", classifyError(resp.StatusCode, body)
	}

	return resp.Body, resp.Header.Get("Content-Type"), nil
}

func classifyError(statusCode int, body []byte) error {
	// Try to extract message from JSON response
	var errResp struct {
		Metadata struct {
			Message string `json:"message"`
		} `json:"metadata"`
	}
	msg := string(body)
	if json.Unmarshal(body, &errResp) == nil && errResp.Metadata.Message != "" {
		msg = errResp.Metadata.Message
	}

	switch statusCode {
	case http.StatusBadRequest:
		return &cerrors.ValidationError{Field: "request", Message: msg}
	case http.StatusForbidden:
		return &cerrors.AuthError{Message: msg}
	case http.StatusNotFound:
		return &cerrors.NotFoundError{Resource: "document", ID: ""}
	default:
		return &cerrors.APIError{StatusCode: statusCode, Message: msg}
	}
}

// MaskAPIKey returns a masked version of the API key for display.
func MaskAPIKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "..." + key[len(key)-4:]
}

// FormatDebug returns a debug string for a URL (with API key masked).
func FormatDebug(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	q := u.Query()
	if q.Get("Subscription-Key") != "" {
		q.Set("Subscription-Key", "****")
		u.RawQuery = q.Encode()
	}
	return fmt.Sprintf("GET %s", u.String())
}
```

- [ ] **Step 5: Run tests to verify they pass**

```bash
go test ./internal/api/... -v
```

Expected: Both tests PASS.

- [ ] **Step 6: Commit**

```bash
git add internal/api/ internal/model/
git commit -m "feat: add API client with error classification and document models"
```

---

## Task 6: API Documents Methods

**Files:**
- Create: `internal/api/documents.go`
- Create: `internal/api/documents_test.go`

- [ ] **Step 1: Write documents API tests**

```go
// internal/api/documents_test.go
package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const sampleDocListResponse = `{
  "metadata": {
    "title": "提出された書類を把握するためのAPI",
    "parameter": {"date": "2024-04-01", "type": "2"},
    "resultset": {"count": 1},
    "processDateTime": "2024-04-01 13:01",
    "status": "200",
    "message": "OK"
  },
  "results": [
    {
      "seqNumber": 1,
      "docID": "S100ABCD",
      "edinetCode": "E10001",
      "secCode": "10000",
      "JCN": "6000012010023",
      "filerName": "テスト株式会社",
      "fundCode": "G00001",
      "ordinanceCode": "030",
      "formCode": "04A000",
      "docTypeCode": "030",
      "periodStart": null,
      "periodEnd": null,
      "submitDateTime": "2024-04-01 12:34",
      "docDescription": "有価証券届出書",
      "issuerEdinetCode": null,
      "subjectEdinetCode": null,
      "subsidiaryEdinetCode": null,
      "currentReportReason": null,
      "parentDocID": null,
      "opeDateTime": null,
      "withdrawalStatus": "0",
      "docInfoEditStatus": "0",
      "disclosureStatus": "0",
      "xbrlFlag": "1",
      "pdfFlag": "1",
      "attachDocFlag": "1",
      "englishDocFlag": "0",
      "csvFlag": "1",
      "legalStatus": "1"
    }
  ]
}`

func TestListDocuments(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/documents.json" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("date") != "2024-04-01" {
			t.Errorf("unexpected date: %s", r.URL.Query().Get("date"))
		}
		if r.URL.Query().Get("type") != "2" {
			t.Errorf("unexpected type: %s", r.URL.Query().Get("type"))
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(sampleDocListResponse))
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, APIKey: "test-key", HTTP: ts.Client()}
	resp, err := c.ListDocuments("2024-04-01", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Metadata.ResultSet.Count != 1 {
		t.Errorf("expected count 1, got %d", resp.Metadata.ResultSet.Count)
	}
	if len(resp.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(resp.Results))
	}
	if resp.Results[0].DocID != "S100ABCD" {
		t.Errorf("expected docID S100ABCD, got %s", resp.Results[0].DocID)
	}
	if resp.Results[0].FilerName != "テスト株式会社" {
		t.Errorf("expected filerName テスト株式会社, got %s", resp.Results[0].FilerName)
	}
}

func TestDownloadDocument(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/documents/S100ABCD" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("type") != "2" {
			t.Errorf("unexpected type: %s", r.URL.Query().Get("type"))
		}
		w.Header().Set("Content-Type", "application/pdf")
		w.Write([]byte("%PDF-fake-content"))
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, APIKey: "test-key", HTTP: ts.Client()}
	body, contentType, err := c.DownloadDocument("S100ABCD", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer body.Close()

	if contentType != "application/pdf" {
		t.Errorf("expected content-type application/pdf, got %s", contentType)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/api/... -v -run TestListDocuments
```

Expected: FAIL — `ListDocuments` not defined.

- [ ] **Step 3: Create documents.go**

```go
// internal/api/documents.go
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/planitaicojp/edinet-cli/internal/model"
)

// ListDocuments calls the document list API.
func (c *Client) ListDocuments(date string, typ int) (*model.DocumentListResponse, error) {
	params := map[string]string{
		"date": date,
		"type": strconv.Itoa(typ),
	}
	body, err := c.Get("/documents.json", params)
	if err != nil {
		return nil, err
	}

	var resp model.DocumentListResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &resp, nil
}

// DownloadDocument calls the document retrieval API and returns the binary body.
func (c *Client) DownloadDocument(docID string, typ int) (io.ReadCloser, string, error) {
	params := map[string]string{
		"type": strconv.Itoa(typ),
	}
	return c.GetBinary(fmt.Sprintf("/documents/%s", docID), params)
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
go test ./internal/api/... -v
```

Expected: All 4 tests PASS (2 from client_test + 2 from documents_test).

- [ ] **Step 5: Commit**

```bash
git add internal/api/documents.go internal/api/documents_test.go
git commit -m "feat: add ListDocuments and DownloadDocument API methods"
```

---

## Task 7: Root Command + Version

**Files:**
- Create: `cmd/root.go`
- Create: `cmd/version.go`

- [ ] **Step 1: Create root.go**

```go
// cmd/root.go
package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "edinet",
	Short: "EDINET API CLI - 金融庁 開示書類の検索・取得ツール",
	Long: `edinet は、金融庁が提供する EDINET API v2 を操作するための CLI ツールです。
有価証券報告書等の開示書類を検索・取得できます。`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().StringP("format", "f", "", "出力形式 (table, json, csv)")
	rootCmd.PersistentFlags().String("api-key", "", "EDINET API キー")
	rootCmd.PersistentFlags().Bool("verbose", false, "デバッグ出力を有効化")
	rootCmd.PersistentFlags().Bool("no-color", false, "カラー出力を無効化")
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
```

- [ ] **Step 2: Create version.go**

```go
// cmd/version.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// version is set via ldflags at build time.
var version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "バージョン情報を表示",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("edinet version %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
```

- [ ] **Step 3: Add cobra dependency and verify build**

```bash
go get github.com/spf13/cobra
go build -o edinet .
./edinet version
```

Expected: `edinet version dev`

- [ ] **Step 4: Verify help**

```bash
./edinet --help
```

Expected: Shows usage with global flags.

- [ ] **Step 5: Commit**

```bash
git add cmd/ main.go go.mod go.sum
git commit -m "feat: add root command with global flags and version subcommand"
```

---

## Task 8: Command Utilities (cmdutil)

**Files:**
- Create: `cmd/cmdutil/client.go`
- Create: `cmd/cmdutil/format.go`
- Create: `cmd/cmdutil/flags.go`

- [ ] **Step 1: Create flags.go**

```go
// cmd/cmdutil/flags.go
package cmdutil

import "github.com/spf13/cobra"

// GetStringFlag returns a string flag value from the command or its parents.
func GetStringFlag(cmd *cobra.Command, name string) string {
	val, _ := cmd.Flags().GetString(name)
	return val
}

// GetBoolFlag returns a bool flag value from the command or its parents.
func GetBoolFlag(cmd *cobra.Command, name string) bool {
	val, _ := cmd.Flags().GetBool(name)
	return val
}
```

- [ ] **Step 2: Create format.go**

```go
// cmd/cmdutil/format.go
package cmdutil

import (
	"os"

	"github.com/planitaicojp/edinet-cli/internal/config"
	"github.com/planitaicojp/edinet-cli/internal/output"
	"github.com/spf13/cobra"
)

// GetFormatter returns the appropriate output formatter based on flag > env > config > default.
func GetFormatter(cmd *cobra.Command) output.Formatter {
	format := GetStringFlag(cmd, "format")
	if format == "" {
		format = os.Getenv(config.EnvFormat)
	}
	if format == "" {
		cfg, err := config.Load()
		if err == nil && cfg.DefaultFormat != "" {
			format = cfg.DefaultFormat
		}
	}
	return output.New(format)
}
```

- [ ] **Step 3: Create client.go**

```go
// cmd/cmdutil/client.go
package cmdutil

import (
	"fmt"

	"github.com/planitaicojp/edinet-cli/internal/api"
	"github.com/planitaicojp/edinet-cli/internal/config"
	cerrors "github.com/planitaicojp/edinet-cli/internal/errors"
	"github.com/spf13/cobra"
)

// NewClient creates an API client from command flags, env vars, and config.
func NewClient(cmd *cobra.Command) (*api.Client, error) {
	apiKey := config.ResolveAPIKey(GetStringFlag(cmd, "api-key"))
	if apiKey == "" {
		return nil, &cerrors.AuthError{
			Message: fmt.Sprintf("API キーが設定されていません。%s 環境変数を設定するか、'edinet config set api-key <key>' を実行してください", config.EnvAPIKey),
		}
	}

	client := api.NewClient(apiKey)
	client.Debug = GetBoolFlag(cmd, "verbose")
	return client, nil
}
```

- [ ] **Step 4: Verify build**

```bash
go build ./...
```

Expected: Success.

- [ ] **Step 5: Commit**

```bash
git add cmd/cmdutil/
git commit -m "feat: add cmdutil helpers (client factory, format resolver, flags)"
```

---

## Task 9: Document Commands

**Files:**
- Create: `cmd/document/document.go`
- Create: `cmd/document/list.go`
- Create: `cmd/document/download.go`

- [ ] **Step 1: Create document.go (group command)**

```go
// cmd/document/document.go
package document

import "github.com/spf13/cobra"

// Cmd is the parent command for document subcommands.
var Cmd = &cobra.Command{
	Use:     "document",
	Aliases: []string{"doc"},
	Short:   "開示書類の一覧取得・ダウンロード",
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(downloadCmd)
}
```

- [ ] **Step 2: Create list.go**

```go
// cmd/document/list.go
package document

import (
	"fmt"
	"os"
	"time"

	"github.com/planitaicojp/edinet-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "指定日付の書類一覧を取得",
	Long:  `指定した日付に提出された開示書類の一覧を取得します。`,
	Example: `  edinet document list --date 2024-04-01
  edinet document list --date today --type 2
  edinet document list --date yesterday --format json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}

		date, _ := cmd.Flags().GetString("date")
		date = resolveDate(date)

		typ, _ := cmd.Flags().GetInt("type")

		resp, err := client.ListDocuments(date, typ)
		if err != nil {
			return err
		}

		formatter := cmdutil.GetFormatter(cmd)

		if typ == 2 && len(resp.Results) > 0 {
			return formatter.Format(os.Stdout, resp.Results)
		}

		return formatter.Format(os.Stdout, resp.Metadata)
	},
}

func init() {
	listCmd.Flags().String("date", "", "対象日付 (YYYY-MM-DD, today, yesterday) [必須]")
	listCmd.Flags().Int("type", 1, "取得情報 (1: メタデータのみ, 2: 提出書類一覧+メタデータ)")
	listCmd.MarkFlagRequired("date")
}

func resolveDate(s string) string {
	now := time.Now()
	switch s {
	case "today":
		return now.Format("2006-01-02")
	case "yesterday":
		return now.AddDate(0, 0, -1).Format("2006-01-02")
	default:
		return s
	}
}

func init() {
}
```

Remove the duplicate empty `init()` at the bottom. The file should end after `resolveDate`.

Actually, let me fix that — the file should be:

```go
// cmd/document/list.go
package document

import (
	"fmt"
	"os"
	"time"

	"github.com/planitaicojp/edinet-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "指定日付の書類一覧を取得",
	Long:  `指定した日付に提出された開示書類の一覧を取得します。`,
	Example: `  edinet document list --date 2024-04-01
  edinet document list --date today --type 2
  edinet document list --date yesterday --format json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}

		date, _ := cmd.Flags().GetString("date")
		date = resolveDate(date)

		typ, _ := cmd.Flags().GetInt("type")

		resp, err := client.ListDocuments(date, typ)
		if err != nil {
			return err
		}

		formatter := cmdutil.GetFormatter(cmd)

		if typ == 2 && len(resp.Results) > 0 {
			return formatter.Format(os.Stdout, resp.Results)
		}

		// type=1: show metadata
		fmt.Fprintf(os.Stdout, "日付: %s / 件数: %d / 更新日時: %s\n",
			resp.Metadata.Parameter.Date,
			resp.Metadata.ResultSet.Count,
			resp.Metadata.ProcessDateTime)
		return nil
	},
}

func init() {
	listCmd.Flags().String("date", "", "対象日付 (YYYY-MM-DD, today, yesterday) [必須]")
	listCmd.Flags().Int("type", 1, "取得情報 (1: メタデータのみ, 2: 提出書類一覧+メタデータ)")
	listCmd.MarkFlagRequired("date")
}

func resolveDate(s string) string {
	now := time.Now()
	switch s {
	case "today":
		return now.Format("2006-01-02")
	case "yesterday":
		return now.AddDate(0, 0, -1).Format("2006-01-02")
	default:
		return s
	}
}
```

- [ ] **Step 3: Create download.go**

```go
// cmd/document/download.go
package document

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/planitaicojp/edinet-cli/cmd/cmdutil"
	cerrors "github.com/planitaicojp/edinet-cli/internal/errors"
	"github.com/spf13/cobra"
)

var docTypeMap = map[string]int{
	"xbrl":    1,
	"pdf":     2,
	"attach":  3,
	"english": 4,
	"csv":     5,
}

var docTypeExt = map[string]string{
	"xbrl":    ".zip",
	"pdf":     ".pdf",
	"attach":  ".zip",
	"english": ".zip",
	"csv":     ".zip",
}

var downloadCmd = &cobra.Command{
	Use:   "download <docID>",
	Short: "書類をダウンロード",
	Long:  `指定した書類管理番号の書類をダウンロードします。`,
	Example: `  edinet document download S100ABCD
  edinet document download S100ABCD --type pdf
  edinet document download S100ABCD --type csv --output ./data/`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cmdutil.NewClient(cmd)
		if err != nil {
			return err
		}

		docID := args[0]
		typeName, _ := cmd.Flags().GetString("type")
		outputDir, _ := cmd.Flags().GetString("output")

		apiType, ok := docTypeMap[typeName]
		if !ok {
			return &cerrors.ValidationError{
				Field:   "type",
				Message: fmt.Sprintf("不正な書類タイプ: %s (xbrl, pdf, attach, english, csv)", typeName),
			}
		}

		body, _, err := client.DownloadDocument(docID, apiType)
		if err != nil {
			return err
		}
		defer body.Close()

		ext := docTypeExt[typeName]
		filename := fmt.Sprintf("%s_%s%s", docID, typeName, ext)
		outPath := filepath.Join(outputDir, filename)

		f, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("ファイル作成エラー: %w", err)
		}
		defer f.Close()

		n, err := io.Copy(f, body)
		if err != nil {
			return fmt.Errorf("ダウンロードエラー: %w", err)
		}

		fmt.Fprintf(os.Stderr, "ダウンロード完了: %s (%d bytes)\n", outPath, n)
		return nil
	},
}

func init() {
	downloadCmd.Flags().StringP("type", "t", "xbrl", "書類タイプ (xbrl, pdf, attach, english, csv)")
	downloadCmd.Flags().StringP("output", "o", ".", "保存先ディレクトリ")
}
```

- [ ] **Step 4: Register document command in root.go**

Add to `cmd/root.go`:

```go
import (
	"github.com/planitaicojp/edinet-cli/cmd/document"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.PersistentFlags().StringP("format", "f", "", "出力形式 (table, json, csv)")
	rootCmd.PersistentFlags().String("api-key", "", "EDINET API キー")
	rootCmd.PersistentFlags().Bool("verbose", false, "デバッグ出力を有効化")
	rootCmd.PersistentFlags().Bool("no-color", false, "カラー出力を無効化")

	rootCmd.AddCommand(document.Cmd)
}
```

- [ ] **Step 5: Build and verify**

```bash
go build -o edinet .
./edinet document --help
./edinet document list --help
./edinet document download --help
```

Expected: Help text for each command.

- [ ] **Step 6: Commit**

```bash
git add cmd/document/ cmd/root.go
git commit -m "feat: add document list and download commands"
```

---

## Task 10: Config Commands

**Files:**
- Create: `cmd/config/config.go`
- Create: `cmd/config/set.go`
- Create: `cmd/config/show.go`

- [ ] **Step 1: Create config.go (group command)**

```go
// cmd/config/config.go
package config

import "github.com/spf13/cobra"

// Cmd is the parent command for config subcommands.
var Cmd = &cobra.Command{
	Use:   "config",
	Short: "設定の表示・変更",
}

func init() {
	Cmd.AddCommand(setCmd)
	Cmd.AddCommand(showCmd)
}
```

- [ ] **Step 2: Create set.go**

```go
// cmd/config/set.go
package config

import (
	"fmt"

	internalConfig "github.com/planitaicojp/edinet-cli/internal/config"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "設定値を変更",
	Example: `  edinet config set api-key YOUR_API_KEY
  edinet config set default-format json`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key, value := args[0], args[1]

		cfg, err := internalConfig.Load()
		if err != nil {
			return err
		}

		switch key {
		case "api-key":
			cfg.APIKey = value
		case "default-format":
			cfg.DefaultFormat = value
		default:
			return fmt.Errorf("不明な設定キー: %s (api-key, default-format)", key)
		}

		if err := internalConfig.Save(cfg); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "設定を保存しました: %s\n", key)
		return nil
	},
}
```

- [ ] **Step 3: Create show.go**

```go
// cmd/config/show.go
package config

import (
	"fmt"

	"github.com/planitaicojp/edinet-cli/internal/api"
	internalConfig "github.com/planitaicojp/edinet-cli/internal/config"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "現在の設定を表示",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := internalConfig.Load()
		if err != nil {
			return err
		}

		w := cmd.OutOrStdout()

		maskedKey := "(未設定)"
		if cfg.APIKey != "" {
			maskedKey = api.MaskAPIKey(cfg.APIKey)
		}

		fmt.Fprintf(w, "config_dir:      %s\n", internalConfig.ConfigDir())
		fmt.Fprintf(w, "api_key:         %s\n", maskedKey)
		fmt.Fprintf(w, "default_format:  %s\n", cfg.DefaultFormat)

		// Show resolved values
		resolved := internalConfig.ResolveAPIKey("")
		if resolved != "" && resolved != cfg.APIKey {
			fmt.Fprintf(w, "\n(環境変数 %s により上書きされています)\n", internalConfig.EnvAPIKey)
		}

		return nil
	},
}
```

- [ ] **Step 4: Register config command in root.go**

Update `cmd/root.go` imports and init:

```go
import (
	configCmd "github.com/planitaicojp/edinet-cli/cmd/config"
	"github.com/planitaicojp/edinet-cli/cmd/document"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.PersistentFlags().StringP("format", "f", "", "出力形式 (table, json, csv)")
	rootCmd.PersistentFlags().String("api-key", "", "EDINET API キー")
	rootCmd.PersistentFlags().Bool("verbose", false, "デバッグ出力を有効化")
	rootCmd.PersistentFlags().Bool("no-color", false, "カラー出力を無効化")

	rootCmd.AddCommand(document.Cmd)
	rootCmd.AddCommand(configCmd.Cmd)
}
```

- [ ] **Step 5: Build and verify**

```bash
go build -o edinet .
./edinet config --help
./edinet config show
./edinet config set api-key test-12345
./edinet config show
```

Expected: config show displays the saved key (masked).

- [ ] **Step 6: Commit**

```bash
git add cmd/config/ cmd/root.go
git commit -m "feat: add config set and show commands"
```

---

## Task 11: Makefile + golangci-lint

**Files:**
- Create: `Makefile`
- Create: `.golangci.yml`

- [ ] **Step 1: Create .golangci.yml**

```yaml
run:
  timeout: 5m

linters:
  enable:
    - errcheck
    - govet
    - staticcheck
    - unused
    - gosimple
    - ineffassign
    - misspell
    - gofmt
    - goimports

linters-settings:
  misspell:
    locale: US

issues:
  exclude-use-default: false
```

- [ ] **Step 2: Create Makefile**

```makefile
BINARY := edinet
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -s -w -X github.com/planitaicojp/edinet-cli/cmd.version=$(VERSION)

.PHONY: build test lint clean install

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

test:
	go test ./... -v

lint:
	golangci-lint run ./...

fmt:
	gofmt -w .
	goimports -w .

clean:
	rm -f $(BINARY)
	rm -rf dist/

install:
	go install -ldflags "$(LDFLAGS)" .

all: lint test build
```

- [ ] **Step 3: Run lint to verify**

```bash
# Install golangci-lint if needed:
# go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
make lint
```

Expected: No lint errors (or fix any found).

- [ ] **Step 4: Run all make targets**

```bash
make test
make build
./edinet version
```

Expected: Tests pass, binary built, version shown.

- [ ] **Step 5: Commit**

```bash
git add Makefile .golangci.yml
git commit -m "feat: add Makefile with build, test, lint targets"
```

---

## Task 12: goreleaser + Homebrew + Scoop

**Files:**
- Create: `.goreleaser.yaml`
- Create: `.gitignore`
- Create: `LICENSE`

- [ ] **Step 1: Create .gitignore**

```
# Binaries
edinet
edinet.exe
dist/

# IDE
.idea/
.vscode/
*.swp

# OS
.DS_Store
Thumbs.db

# Test
coverage.out
```

- [ ] **Step 2: Create LICENSE (Apache-2.0)**

```
                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   ... (standard Apache 2.0 text)

   Copyright 2026 PlanitAI Co., Ltd.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
```

Use the full standard Apache 2.0 license text.

- [ ] **Step 3: Create .goreleaser.yaml**

```yaml
version: 2

builds:
  - binary: edinet
    ldflags:
      - -s -w -X github.com/planitaicojp/edinet-cli/cmd.version={{.Version}}
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64

archives:
  - format: tar.gz
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
    format_overrides:
      - goos: windows
        format: zip

checksum:
  name_template: "checksums.txt"

brews:
  - repository:
      owner: planitaicojp
      name: homebrew-tap
      token: "{{ .Env.HOMEBREW_TAP_TOKEN }}"
    name: edinet
    homepage: "https://github.com/planitaicojp/edinet-cli"
    description: "CLI tool for EDINET API v2 - 金融庁 開示書類の検索・取得"
    license: "Apache-2.0"

scoops:
  - repository:
      owner: planitaicojp
      name: bucket
      token: "{{ .Env.SCOOP_BUCKET_TOKEN }}"
    name: edinet
    homepage: "https://github.com/planitaicojp/edinet-cli"
    description: "CLI tool for EDINET API v2 - 金融庁 開示書類の検索・取得"
    license: "Apache-2.0"

changelog:
  sort: asc
  filters:
    exclude:
      - "^docs:"
      - "^test:"
      - "^chore:"
```

- [ ] **Step 4: Verify goreleaser config**

```bash
# Install goreleaser if needed, then check:
goreleaser check
```

If goreleaser is not installed, verify YAML syntax manually. The config will be validated on first release.

- [ ] **Step 5: Test local build with goreleaser**

```bash
goreleaser build --snapshot --clean
```

Expected: Binaries created in `dist/` for all OS/arch combinations.

- [ ] **Step 6: Commit**

```bash
git add .goreleaser.yaml .gitignore LICENSE
git commit -m "feat: add goreleaser config with Homebrew and Scoop support"
```

---

## Task 13: README.md

**Files:**
- Create: `README.md`

- [ ] **Step 1: Create README.md**

```markdown
# edinet-cli

金融庁が提供する [EDINET API v2](https://disclosure2dl.edinet-fsa.go.jp/guide/static/disclosure/WZEK0110.html) を操作するための CLI ツールです。
有価証券報告書等の開示書類を検索・取得できます。

## インストール

### Homebrew (macOS / Linux)

```bash
brew install planitaicojp/tap/edinet
```

### Scoop (Windows)

```powershell
scoop bucket add planitaicojp https://github.com/planitaicojp/bucket
scoop install edinet
```

### ソースからビルド

```bash
go install github.com/planitaicojp/edinet-cli@latest
```

### リリースバイナリ

[Releases](https://github.com/planitaicojp/edinet-cli/releases) ページからダウンロード：

**Linux (amd64)**

```bash
VERSION=$(curl -s https://api.github.com/repos/planitaicojp/edinet-cli/releases/latest | grep tag_name | cut -d '"' -f4)
curl -Lo edinet.tar.gz "https://github.com/planitaicojp/edinet-cli/releases/download/${VERSION}/edinet-cli_${VERSION#v}_linux_amd64.tar.gz"
tar xzf edinet.tar.gz edinet
sudo mv edinet /usr/local/bin/
rm edinet.tar.gz
```

**macOS (Apple Silicon)**

```bash
VERSION=$(curl -s https://api.github.com/repos/planitaicojp/edinet-cli/releases/latest | grep tag_name | cut -d '"' -f4)
curl -Lo edinet.tar.gz "https://github.com/planitaicojp/edinet-cli/releases/download/${VERSION}/edinet-cli_${VERSION#v}_darwin_arm64.tar.gz"
tar xzf edinet.tar.gz edinet
sudo mv edinet /usr/local/bin/
rm edinet.tar.gz
```

**Windows (amd64)**

```powershell
$version = (Invoke-RestMethod https://api.github.com/repos/planitaicojp/edinet-cli/releases/latest).tag_name
$v = $version -replace '^v', ''
Invoke-WebRequest -Uri "https://github.com/planitaicojp/edinet-cli/releases/download/$version/edinet-cli_${v}_windows_amd64.zip" -OutFile edinet.zip
Expand-Archive edinet.zip -DestinationPath .
Remove-Item edinet.zip
```

## セットアップ

### API キーの取得

1. [EDINET API](https://api.edinet-fsa.go.jp/api/auth/index.aspx?mode=1) でアカウントを作成
2. API キーを発行

### API キーの設定

```bash
# 方法1: 環境変数 (推奨)
export EDINET_API_KEY="your-api-key"

# 方法2: config ファイル
edinet config set api-key "your-api-key"
```

## 使い方

### 書類一覧の取得

```bash
# メタデータの取得
edinet document list --date 2024-04-01

# 提出書類一覧の取得
edinet document list --date today --type 2

# JSON形式で出力
edinet document list --date yesterday --type 2 --format json
```

### 書類のダウンロード

```bash
# XBRL (デフォルト)
edinet document download S100ABCD

# PDF
edinet document download S100ABCD --type pdf

# CSV
edinet document download S100ABCD --type csv --output ./data/
```

### 設定の確認

```bash
edinet config show
```

## 環境変数

| 変数名 | 説明 | デフォルト |
|--------|------|-----------|
| `EDINET_API_KEY` | API キー | - |
| `EDINET_FORMAT` | 出力形式 (table, json, csv) | `table` |
| `EDINET_CONFIG_DIR` | 設定ディレクトリ | `~/.config/edinet` |
| `EDINET_DEBUG` | デバッグログ | `false` |

## ライセンス

[Apache License 2.0](LICENSE)
```

- [ ] **Step 2: Commit**

```bash
git add README.md
git commit -m "docs: add README with installation and usage instructions"
```

---

## Task 14: Final Integration Test

- [ ] **Step 1: Clean build**

```bash
rm -f edinet
go build -o edinet .
```

- [ ] **Step 2: Verify all commands**

```bash
./edinet version
./edinet --help
./edinet document --help
./edinet document list --help
./edinet document download --help
./edinet config --help
./edinet config show
```

Expected: All commands show help/output without errors.

- [ ] **Step 3: Run all tests**

```bash
go test ./... -v
```

Expected: All tests pass.

- [ ] **Step 4: Verify goreleaser snapshot**

```bash
goreleaser build --snapshot --clean
ls dist/
```

Expected: Binaries for linux/darwin/windows, amd64/arm64.

- [ ] **Step 5: Final commit (if any changes needed)**

```bash
git add -A
git commit -m "chore: final integration verification"
```
