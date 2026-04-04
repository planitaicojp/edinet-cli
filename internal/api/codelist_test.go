package api

import (
	"strings"
	"testing"

	"github.com/planitaicojp/edinet-cli/internal/model"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func TestParseCodeListCSV(t *testing.T) {
	// Create Shift_JIS encoded CSV content
	csvContent := "EDINETコードリスト\r\n" +
		"EDINETコード,提出者種別,上場区分,連結の有無,資本金,決算日,提出者名,提出者名（英字）,提出者名（ヨミ）,所在地,提出者業種,証券コード,提出者法人番号\r\n" +
		"E00001,内国法人・組合,上場,有,100000000,3月31日,テスト株式会社,Test Corp,テストカブシキガイシャ,東京都千代田区,電気機器,1234,1234567890123\r\n" +
		"E00002,内国法人・組合,非上場,無,50000000,12月31日,サンプル株式会社,Sample Inc,サンプルカブシキガイシャ,大阪府大阪市,情報・通信業,5678,9876543210123\r\n"

	// Encode to Shift_JIS
	var buf strings.Builder
	w := transform.NewWriter(&buf, japanese.ShiftJIS.NewEncoder())
	_, err := w.Write([]byte(csvContent))
	if err != nil {
		t.Fatalf("failed to encode test data: %v", err)
	}
	_ = w.Close()

	companies, err := ParseCodeListCSV(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(companies) != 2 {
		t.Fatalf("expected 2 companies, got %d", len(companies))
	}

	if companies[0].EdinetCode != "E00001" {
		t.Errorf("expected EdinetCode 'E00001', got %q", companies[0].EdinetCode)
	}
	if companies[0].FilerName != "テスト株式会社" {
		t.Errorf("expected FilerName 'テスト株式会社', got %q", companies[0].FilerName)
	}
	if companies[0].SecCode != "1234" {
		t.Errorf("expected SecCode '1234', got %q", companies[0].SecCode)
	}
	if companies[1].FilerNameEN != "Sample Inc" {
		t.Errorf("expected FilerNameEN 'Sample Inc', got %q", companies[1].FilerNameEN)
	}
}

func TestSearchCompanies(t *testing.T) {
	companies := []model.Company{
		{EdinetCode: "E00001", FilerName: "トヨタ自動車", SecCode: "7203", FilerNameEN: "Toyota Motor"},
		{EdinetCode: "E00002", FilerName: "ソニーグループ", SecCode: "6758", FilerNameEN: "Sony Group"},
		{EdinetCode: "E00003", FilerName: "任天堂", SecCode: "7974", FilerNameEN: "Nintendo"},
	}

	tests := []struct {
		query    string
		by       string
		expected int
	}{
		{"トヨタ", "all", 1},
		{"toyota", "name", 1},
		{"7203", "code", 1},
		{"E00002", "edinet-code", 1},
		{"グループ", "all", 1},
		{"missing", "all", 0},
	}

	for _, tt := range tests {
		results := SearchCompanies(companies, tt.query, tt.by)
		if len(results) != tt.expected {
			t.Errorf("SearchCompanies(%q, %q): got %d results, want %d", tt.query, tt.by, len(results), tt.expected)
		}
	}
}
