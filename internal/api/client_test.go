package api

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestClient_SubscriptionKey(t *testing.T) {
	var gotKey string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotKey = r.URL.Query().Get("Subscription-Key")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"metadata":{"title":"test","parameter":{"date":"2024-01-01","type":"1"},"resultset":{"count":0},"processDateTime":"2024-01-01 00:00","status":"200","message":"OK"}}`))
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
		_, _ = w.Write([]byte(`{"metadata":{"status":"403","message":"Invalid API Key"}}`))
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

func TestClient_RetryOn500(t *testing.T) {
	var attempts int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&attempts, 1)
		if n <= 2 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"metadata":{"message":"server error"}}`))
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"metadata":{"title":"test","parameter":{"date":"2024-01-01","type":"1"},"resultset":{"count":0},"processDateTime":"2024-01-01 00:00","status":"200","message":"OK"}}`))
	}))
	defer ts.Close()

	c := &Client{
		BaseURL:     ts.URL,
		APIKey:      "test-key",
		HTTP:        ts.Client(),
		BackoffFunc: func(int) time.Duration { return 0 },
	}

	_, err := c.Get("/documents.json", map[string]string{"date": "2024-01-01"})
	if err != nil {
		t.Fatalf("expected success after retry, got: %v", err)
	}
	if atomic.LoadInt32(&attempts) != 3 {
		t.Errorf("expected 3 attempts, got %d", atomic.LoadInt32(&attempts))
	}
}

func TestClient_RetryExhausted(t *testing.T) {
	var attempts int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"metadata":{"message":"server error"}}`))
	}))
	defer ts.Close()

	c := &Client{
		BaseURL:     ts.URL,
		APIKey:      "test-key",
		HTTP:        ts.Client(),
		BackoffFunc: func(int) time.Duration { return 0 },
	}

	_, err := c.Get("/documents.json", map[string]string{"date": "2024-01-01"})
	if err == nil {
		t.Fatal("expected error after exhausted retries")
	}
	// initial + 3 retries = 4
	if atomic.LoadInt32(&attempts) != 4 {
		t.Errorf("expected 4 attempts, got %d", atomic.LoadInt32(&attempts))
	}
}

func TestMaskAPIKey(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"short", "****"},
		{"12345678abcdefgh", "1234...efgh"},
	}
	for _, tt := range tests {
		got := MaskAPIKey(tt.input)
		if got != tt.expected {
			t.Errorf("MaskAPIKey(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestMaskSubscriptionKey(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"https://api.example.com?Subscription-Key=secret123&date=2024-01-01",
			"https://api.example.com?Subscription-Key=****&date=2024-01-01",
		},
		{
			"https://api.example.com?date=2024-01-01&Subscription-Key=secret123",
			"https://api.example.com?date=2024-01-01&Subscription-Key=****",
		},
		{
			"https://api.example.com?date=2024-01-01",
			"https://api.example.com?date=2024-01-01",
		},
	}
	for _, tt := range tests {
		got := maskSubscriptionKey(tt.input)
		if got != tt.expected {
			t.Errorf("maskSubscriptionKey(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}
