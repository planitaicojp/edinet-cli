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
