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
		_, _ = w.Write([]byte(sampleDocListResponse))
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
		_, _ = w.Write([]byte("%PDF-fake-content"))
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, APIKey: "test-key", HTTP: ts.Client()}
	body, contentType, err := c.DownloadDocument("S100ABCD", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = body.Close() }()

	if contentType != "application/pdf" {
		t.Errorf("expected content-type application/pdf, got %s", contentType)
	}
}
