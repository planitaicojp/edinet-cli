package api

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/planitaicojp/edinet-cli/internal/config"
	"github.com/planitaicojp/edinet-cli/internal/model"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const codeListURL = "https://disclosure2dl.edinet-fsa.go.jp/searchdocument/codelist/Edinetcode.zip"

// DownloadCodeList downloads the EDINET code list ZIP and returns the response body.
func (c *Client) DownloadCodeList() (io.ReadCloser, error) {
	resp, err := c.doGet(codeListURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download code list: %w", err)
	}
	if resp.StatusCode != 200 {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("failed to download code list: HTTP %d", resp.StatusCode)
	}
	return resp.Body, nil
}

// CacheFilePath returns the path to the cached code list CSV.
func CacheFilePath() string {
	return filepath.Join(config.ConfigDir(), "codelist.csv")
}

// LoadCodeList reads the cached CSV file and returns Company records.
func LoadCodeList(path string) ([]model.Company, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("コードリストが見つかりません。'edinet company list --update' で取得してください: %w", err)
	}
	defer func() { _ = f.Close() }()

	return ParseCodeListCSV(f)
}

// ParseCodeListCSV parses EDINET code list CSV (Shift_JIS encoded, header rows to skip).
func ParseCodeListCSV(r io.Reader) ([]model.Company, error) {
	// EDINET CSV is Shift_JIS encoded
	sjisReader := transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
	cr := csv.NewReader(sjisReader)
	cr.LazyQuotes = true
	cr.FieldsPerRecord = -1

	// Skip header row (first line is a title row, second is column headers)
	if _, err := cr.Read(); err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}
	if _, err := cr.Read(); err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	var companies []model.Company
	for {
		record, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV row: %w", err)
		}
		if len(record) < 13 {
			continue
		}
		companies = append(companies, model.Company{
			EdinetCode:    strings.TrimSpace(record[0]),
			FilerType:     strings.TrimSpace(record[1]),
			ListingStatus: strings.TrimSpace(record[2]),
			Consolidated:  strings.TrimSpace(record[3]),
			Capital:       strings.TrimSpace(record[4]),
			FiscalYearEnd: strings.TrimSpace(record[5]),
			FilerName:     strings.TrimSpace(record[6]),
			FilerNameEN:   strings.TrimSpace(record[7]),
			FilerNameKana: strings.TrimSpace(record[8]),
			Address:       strings.TrimSpace(record[9]),
			Industry:      strings.TrimSpace(record[10]),
			SecCode:       strings.TrimSpace(record[11]),
			JCN:           strings.TrimSpace(record[12]),
		})
	}

	return companies, nil
}

// SearchCompanies filters companies by query and field.
func SearchCompanies(companies []model.Company, query, by string) []model.Company {
	query = strings.ToLower(query)
	var results []model.Company
	for _, c := range companies {
		if matchCompany(c, query, by) {
			results = append(results, c)
		}
	}
	return results
}

func matchCompany(c model.Company, query, by string) bool {
	switch by {
	case "name":
		return containsLower(c.FilerName, query) || containsLower(c.FilerNameEN, query) || containsLower(c.FilerNameKana, query)
	case "code":
		return containsLower(c.SecCode, query)
	case "edinet-code":
		return containsLower(c.EdinetCode, query)
	default: // "all"
		return containsLower(c.EdinetCode, query) ||
			containsLower(c.SecCode, query) ||
			containsLower(c.JCN, query) ||
			containsLower(c.FilerName, query) ||
			containsLower(c.FilerNameEN, query) ||
			containsLower(c.FilerNameKana, query)
	}
}

func containsLower(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), substr)
}
